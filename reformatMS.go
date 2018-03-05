package main

import (
	"encoding/csv"
	"log"
	"os"
	"reformatMS/fileHandler"
	"reformatMS/input"
	"strconv"
	"fmt"
)

func main() {
	fmt.Println("RULES: sample number must be even (treatment and control.")
	fmt.Println("Filenames must be entered within ' ' and must end in csv.")
	fmt.Println("The SWATH-MS file copied from the PeakView .xslx output file must be saved as .csv - only the ion sheet.")
	fmt.Println("Biological Replicates should be the name in the intensity column (name of sample) along with _1 if it's the first bioreplicate.")
	fmt.Println("Control should be just the name of the sample, like bio replicate but with the _1.")
	fmt.Println("The FDR file copied from the PeakView .xslx output file must be saved as .csv - only the FDR sheet.")
	openSWATHfile, err := input.Input("What SWATH-MS file are you opening (written like: SWATH.csv): ")
	if err != nil {
		log.Fatalln(err)
	}
	openFDRfile, err := input.Input("What FDR file are you opening (written like: FDR.csv): ")
	if err != nil {
		log.Fatalln(err)
	}
	filename, err := input.Input("What would you like to name the output file (written like: MSstats.csv): ")
	if err != nil {
		log.Fatalln(err)
	}
	swathFile := fileHandler.ReadFile(openSWATHfile, 1)
	samples := len(swathFile.Header) - 9
	fdrFile := fileHandler.ReadFile(openFDRfile, 1)

	fdrMap := make(map[string][]float64)
	for c := range fdrFile.OutputChan {
		fdrFail := 0
		var fdrArray []float64
		for i := 0; i < samples; i++ {
			val, err := strconv.ParseFloat(c[7+i], 64)
			if err != nil {
				log.Fatalln(err)
			}

			if val >= 0.01 {
				fdrFail++
			}
			fdrArray = append(fdrArray, val)
		}
		if fdrFail < samples {

			fdrMap[c[1]] = fdrArray
		}
	}

	o, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	writer := csv.NewWriter(o)
	writer.Comma = ','
	writer.Write([]string{"ProteinName", "PeptideSequence", "PrecursorCharge", "FragmentIon", "ProductCharge", "IsotopeLabelType", "Condition", "BioReplicate", "Run", "Intensity"})
	//log.Println(fdrMap)
	for c := range swathFile.OutputChan {
		if val, ok := fdrMap[c[1]]; ok {

			for i := 0; i < samples; i++ {
				//log.Println(swathFile.Header[9+i])
				row := []string{c[0], c[1], c[3], c[7] + c[8], c[6], "L",
					swathFile.Header[9+i][:len(swathFile.Header[9+i])-2],
					swathFile.Header[9+i],
					strconv.Itoa(i + 1), ""}
				if val[i] < 0.01 {
					row[9] = c[9+i]
				}
				writer.Write(row)
			}
		}
	}
	writer.Flush()
	o.Close()
}
