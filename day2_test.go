package adventofcode2017_test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SpreadsheetRow struct {
	cells []int
}

type Spreadsheet struct {
	rows []SpreadsheetRow
}

func NewSpreadsheetRow(descriptor string) *SpreadsheetRow {
	cell_descriptors := strings.Split(descriptor, "\t")
	row := make([]int, len(cell_descriptors))
	for j := 0; j < len(cell_descriptors); j++ {
		cell_value, err := strconv.Atoi(cell_descriptors[j])
		if err != nil {
			panic(fmt.Sprintf("cannot parse '%s' as an int", cell_descriptors[j]))
		}
		row[j] = cell_value
	}
	return &SpreadsheetRow{row}
}

func NewSpreadsheet(descriptor string) *Spreadsheet {
	row_descriptors := strings.Split(descriptor, "\n")
	rows := make([]SpreadsheetRow, len(row_descriptors))
	for j := 0; j < len(row_descriptors); j++ {
		if len(row_descriptors[j]) > 0 {
			rows[j] = *NewSpreadsheetRow(row_descriptors[j])
		}
	}
	return &Spreadsheet{rows}
}

func (ssr SpreadsheetRow) checksum() int {
	max := ssr.cells[0]
	min := ssr.cells[0]
	for jcell := 1; jcell < len(ssr.cells); jcell++ {
		current := ssr.cells[jcell]
		if current > max {
			max = current
		}
		if current < min {
			min = current
		}
	}
	return max - min
}

func (ss Spreadsheet) checksum() int {
	checksum := 0
	for jrow := 0; jrow < len(ss.rows); jrow++ {
		checksum += ss.rows[jrow].checksum()
	}
	return checksum
}

func (ssr SpreadsheetRow) checksum2() int {
	// find the first two evenly-divisible numbers
	// and return the quotient
	for jcell := 0; jcell < len(ssr.cells)-1; jcell++ {
		jcurr := ssr.cells[jcell]
		for kcell := jcell + 1; kcell < len(ssr.cells); kcell++ {
			kcurr := ssr.cells[kcell]
			if kcurr%jcurr == 0 {
				return kcurr / jcurr
			} else if jcurr%kcurr == 0 {
				return jcurr / kcurr
			}
		}
	}
	return 0
}

func (ss Spreadsheet) checksum2() int {
	checksum := 0
	for jrow := 0; jrow < len(ss.rows); jrow++ {
		checksum += ss.rows[jrow].checksum2()
	}
	return checksum
}

var _ = Describe("Day2", func() {
	Describe("SpreadsheetRow", func() {
		Describe("checksum", func() {
			It("returns the diff between largest and smallest values", func() {
				Expect(SpreadsheetRow{[]int{5, 1, 9, 5}}.checksum()).To(Equal(8))
				Expect(SpreadsheetRow{[]int{7, 5, 3}}.checksum()).To(Equal(4))
				Expect(SpreadsheetRow{[]int{2, 4, 6, 8}}.checksum()).To(Equal(6))
			})
		})

		Describe("checksum2", func() {
			It("returns the quotient of the divisible numbers", func() {
				Expect(SpreadsheetRow{[]int{5, 9, 2, 8}}.checksum2()).To(Equal(4))
				Expect(SpreadsheetRow{[]int{9, 4, 7, 3}}.checksum2()).To(Equal(3))
				Expect(SpreadsheetRow{[]int{3, 8, 6, 5}}.checksum2()).To(Equal(2))
			})
		})
	})
	Describe("Spreadsheet", func() {
		Describe("checksum", func() {
			It("returns the sum of the row checksums", func() {
				rawData := heredoc.Doc(`
					5	1	9	5
					7	5	3
					2	4	6	8`)
				Expect(NewSpreadsheet(rawData).checksum()).To(Equal(18))
			})
		})

		Describe("checksum2", func() {
			It("returns the sum of the row checksums", func() {
				rawData := heredoc.Doc(`
					5	9	2	8
					9	4	7	3
					3	8	6	5`)
				Expect(NewSpreadsheet(rawData).checksum2()).To(Equal(9))
			})
		})
	})

	Describe("puzzle", func() {
		rawData := heredoc.Doc(`
		  409	194	207	470	178	454	235	333	511	103	474	293	525	372	408	428
		  4321	2786	6683	3921	265	262	6206	2207	5712	214	6750	2742	777	5297	3764	167
		  3536	2675	1298	1069	175	145	706	2614	4067	4377	146	134	1930	3850	213	4151
		  2169	1050	3705	2424	614	3253	222	3287	3340	2637	61	216	2894	247	3905	214
		  99	797	80	683	789	92	736	318	103	153	749	631	626	367	110	805
		  2922	1764	178	3420	3246	3456	73	2668	3518	1524	273	2237	228	1826	182	2312
		  2304	2058	286	2258	1607	2492	2479	164	171	663	62	144	1195	116	2172	1839
		  114	170	82	50	158	111	165	164	106	70	178	87	182	101	86	168
		  121	110	51	122	92	146	13	53	34	112	44	160	56	93	82	98
		  4682	642	397	5208	136	4766	180	1673	1263	4757	4680	141	4430	1098	188	1451
		  158	712	1382	170	550	913	191	163	459	1197	1488	1337	900	1182	1018	337
		  4232	236	3835	3847	3881	4180	4204	4030	220	1268	251	4739	246	3798	1885	3244
		  169	1928	3305	167	194	3080	2164	192	3073	1848	426	2270	3572	3456	217	3269
		  140	1005	2063	3048	3742	3361	117	93	2695	1529	120	3480	3061	150	3383	190
		  489	732	57	75	61	797	266	593	324	475	733	737	113	68	267	141
		  3858	202	1141	3458	2507	239	199	4400	3713	3980	4170	227	3968	1688	4352	4168`)
		It("star 1", func() {
			fmt.Println("d2s1: ", NewSpreadsheet(rawData).checksum())
		})
		It("star 2", func() {
			fmt.Println("d2s2: ", NewSpreadsheet(rawData).checksum2())
		})
	})
})
