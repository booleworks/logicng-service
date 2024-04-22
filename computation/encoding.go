package computation

import (
	"fmt"
	"net/http"

	"github.com/booleworks/logicng-go/encoding"
	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-service/config"
	"github.com/booleworks/logicng-service/sio"
)

func HandleEncoding(cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := r.PathValue("enc")
		switch enc {
		case "cc":
			handleEncodingCC(w, r)
		case "pbc":
			handleEncodingPBC(w, r)
		default:
			sio.WriteError(w, r, sio.ErrUnknownPath(r.URL.Path))
		}
	})
}

// @Summary      Encodes a cardinality constraint to a CNF
// @Tags         Encoding
// @Param        algorithm query string false "Encoding algorithm" Enums(pure, ladder, bimander, commander, nested, binary, product, totalizer, mod_totalizer, cardinality_network)
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.FormulaResult
// @Router       /encoding/cc [post]
func handleEncodingCC(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, form formula.Formula) (formula.Formula, sio.ServiceError) {
		if form.Sort() != formula.SortCC {
			return 0, sio.ErrIllegalInput(fmt.Errorf("input is not a cardinality constraint"))
		}
		encCfg, err := extractEncConfig(r)
		if err != nil {
			return 0, sio.ErrIllegalInput(err)
		}
		enc, err := encoding.EncodeCC(fac, form, encCfg)
		if err != nil {
			return 0, sio.ErrIllegalInput(err)
		}
		return fac.And(enc...), nil
	})
}

// @Summary      Encodes a pseudo-Boolean constraint to a CNF
// @Tags         Encoding
// @Param        algorithm query string false "Encoding algorithm" Enums(swc, binary_merge, adder_networks)
// @Param        request body	sio.FormulaInput true "Input Formula"
// @Success      200  {object}  sio.FormulaResult
// @Router       /encoding/pbc [post]
func handleEncodingPBC(w http.ResponseWriter, r *http.Request) {
	transform(w, r, func(fac formula.Factory, form formula.Formula) (formula.Formula, sio.ServiceError) {
		if form.Sort() != formula.SortPBC {
			return 0, sio.ErrIllegalInput(fmt.Errorf("input is not a pseudo-Boolean constraint"))
		}
		encCfg, err := extractEncConfig(r)
		if err != nil {
			return 0, sio.ErrIllegalInput(err)
		}
		enc, err := encoding.EncodePBC(fac, form, encCfg)
		if err != nil {
			return 0, sio.ErrIllegalInput(err)
		}
		return fac.And(enc...), nil
	})
}

func extractEncConfig(r *http.Request) (*encoding.Config, error) {
	encCfg := encoding.DefaultConfig()
	switch algorithm := r.URL.Query().Get("algorithm"); algorithm {
	case "":
		// do nothing
	case "pure":
		encCfg.AMOEncoder = encoding.AMOPure
	case "ladder":
		encCfg.AMOEncoder = encoding.AMOLadder
	case "bimander":
		encCfg.AMOEncoder = encoding.AMOBimander
	case "commander":
		encCfg.AMOEncoder = encoding.AMOCommander
	case "nested":
		encCfg.AMOEncoder = encoding.AMONested
	case "binary":
		encCfg.AMOEncoder = encoding.AMOBinary
	case "product":
		encCfg.AMOEncoder = encoding.AMOProduct
	case "totalizer":
		encCfg.ALKEncoder = encoding.ALKTotalizer
		encCfg.AMKEncoder = encoding.AMKTotalizer
		encCfg.EXKEncoder = encoding.EXKTotalizer
	case "mod_totalizer":
		encCfg.ALKEncoder = encoding.ALKModularTotalizer
		encCfg.AMKEncoder = encoding.AMKModularTotalizer
	case "cardinality_network":
		encCfg.ALKEncoder = encoding.ALKCardinalityNetwork
		encCfg.AMKEncoder = encoding.AMKCardinalityNetwork
		encCfg.EXKEncoder = encoding.EXKCardinalityNetwork
	case "swc":
		encCfg.PBCEncoder = encoding.PBCSWC
	case "binary_merge":
		encCfg.PBCEncoder = encoding.PBCBinaryMerge
	case "adder_networks":
		encCfg.PBCEncoder = encoding.PBCAdderNetworks
	default:
		return nil, fmt.Errorf("unknown encoding algorithm '%s'", algorithm)
	}
	return encCfg, nil
}
