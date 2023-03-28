package main

import (
	"flag"
	"log"
)

var zkBackend = flag.String("backend", "groth16", "Backend used in the ZkSnark circuit ('groth16' or 'plonk')")
var ccsOutput = flag.String("ccs", "./artifacts/zkcensus.ccs", "Output file to the encoded output Gnark Circuit Constrain System, result of circuit compilation")
var srsOutput = flag.String("srs", "./artifacts/zkcensus.srs", "Output file to the encoded output Gnark KZG polynomial commitment, result of circuit compilation")

func main() {
	flag.Parse()

	switch *zkBackend {
	case "plonk":
		ccs, srs, err := compilePlonk()
		if err != nil {
			log.Fatalln(err)
		}
		if err := savePlonk(ccs, srs, *ccsOutput, *srsOutput); err != nil {
			return
		}
	case "groth16":
		ccs, srs, err := compileGroth16()
		if err != nil {
			log.Fatalln(err)
		}
		if err := saveGroth16(ccs, srs, *ccsOutput, *srsOutput); err != nil {
			return
		}
	}

}
