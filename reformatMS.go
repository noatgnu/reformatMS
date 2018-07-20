package main

import (
	"encoding/csv"
	"log"
	"os"
	"github.com/noatgnu/reformatMS/fileHandler"
	"github.com/noatgnu/reformatMS/input"
	"strconv"
		"strings"
	"flag"
)
var swath = flag.String("swath", "", "SWATH File")
var fdr = flag.String("fdr", "", "FDR File")
var out = flag.String("out", "", "Output File")
func init() {
	flag.Parse()
}

func main() {
	log.Println("RULES:\n " +
		"Sample number must be even (treatment and control.\n" +
		"Filenames must be entered within ' ' and must end in csv.\n" +
		"The SWATH-MS file copied from the PeakView .xslx output file must be saved as .csv - only the ion sheet.\n" +
		"Biological Replicates should be the name in the intensity column (name of sample) along with _1 if it's the first bioreplicate.\n" +
		"Control should be just the name of the sample, like bio replicate but with the _1.\n" +
		"The FDR file copied from the PeakView .xslx output file must be saved as .csv - only the FDR sheet.")

	var openSWATHfile, openFDRfile, filename string
	var err error
	if *swath == "" {
		openSWATHfile, err = input.Input("What SWATH-MS file are you opening (written like: SWATH.csv): ")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		openSWATHfile = *swath
	}
	openSWATHfile = input.Clean(openSWATHfile)
	if *fdr == "" {
		openFDRfile, err = input.Input("What FDR file are you opening (written like: FDR.csv): ")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		openFDRfile = *fdr
	}
	openFDRfile = input.Clean(openFDRfile)
	if *out == "" {
		filename, err = input.Input("What would you like to name the output file (written like: MSstats.csv): ")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		filename = *out
	}
	filename = input.Clean(filename)
	log.Printf("Input:\n- SWATH File: %s\n- FDR File: %s\n- Output File: %s ", openSWATHfile, openFDRfile, filename)

	swathFile := fileHandler.ReadFile(openSWATHfile, 1)
	samples := len(swathFile.Header) - 9
	log.Printf("%d Samples", samples)
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
				n := strings.Split(swathFile.Header[9+i], "_")
				row := []string{c[0], c[1], c[3], c[7] + c[8], c[6], "L",
					n[0],
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
	log.Println("Completed.")
}
