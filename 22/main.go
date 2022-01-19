package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Spread struct {
	low  int
	high int
}

func (s Spread) Contains(v int) bool {
	return s.low <= v && v <= s.high
}

func (s Spread) Surrounds(ns Spread) bool {
	return ns.low <= s.low && ns.high >= s.high
}

type Action string

const (
	ON  Action = "on"
	OFF Action = "off"
)

type Command struct {
	action  Action
	x, y, z Spread
}

type RangeMap struct {
	ranges []Spread
	values []bool
}

func NewRangeMap() *RangeMap {
	return &RangeMap{}
}

func (r *RangeMap) Add(newSpread Spread, v bool) {
	fmt.Printf("Adding %+v to %+v\n", newSpread, r.ranges)
	defer func() {
		fmt.Printf("New ranges: %+v\n", r.ranges)
	}()
	if len(r.ranges) == 0 {
		r.ranges = append(r.ranges, newSpread)
		r.values = append(r.values, v)
		return
	}
	idxLow := sort.Search(len(r.ranges), func(i int) bool {
		contains := r.ranges[i].Contains(newSpread.low)
		return contains || r.ranges[i].low >= newSpread.low
	})
	idxHigh := sort.Search(len(r.ranges), func(i int) bool {
		contains := r.ranges[i].Contains(newSpread.high)
		return contains || r.ranges[i].high >= newSpread.high
	})
	fmt.Printf("idxLow, High: %d, %d\n", idxLow, idxHigh)

	if idxLow >= len(r.ranges) {
		r.ranges = append(r.ranges, newSpread)
		r.values = append(r.values, v)
		return
	}

	if idxHigh == 0 {
		if !r.ranges[0].Contains(newSpread.high) {
			r.ranges = append([]Spread{newSpread}, r.ranges...)
			r.values = append([]bool{v}, r.values...)
			return
		}
	}

	if idxLow == idxHigh {
		spreads := append([]Spread{}, r.ranges[:idxLow]...)
		values := append([]bool{}, r.values[:idxLow]...)

		spreads = append(spreads, Spread{r.ranges[idxLow].low, newSpread.low - 1})
		values = append(values, r.values[idxLow])

		spreads = append(spreads, newSpread)
		values = append(values, v)

		spreads = append(spreads, Spread{newSpread.high + 1, r.ranges[idxLow].high})
		values = append(values, r.values[idxLow])

		spreads = append(spreads, r.ranges[idxLow+1:]...)
		values = append(values, r.values[idxLow+1:]...)

		r.ranges = spreads
		r.values = values
		return
	}

	// idxLow != idxHigh
	spreads := append([]Spread{}, r.ranges[:idxLow]...)
	values := append([]bool{}, r.values[:idxLow]...)

	spreads = append(spreads, Spread{r.ranges[idxLow].low, newSpread.low - 1})
	values = append(values, r.values[idxLow])

	spreads = append(spreads, newSpread)
	values = append(values, v)

	if idxHigh < len(r.ranges) {
		spreads = append(spreads, Spread{newSpread.high + 1, r.ranges[idxHigh].high})
		values = append(values, r.values[idxHigh])

		spreads = append(spreads, r.ranges[idxHigh+1:]...)
		values = append(values, r.values[idxHigh+1:]...)
	}

	r.ranges = spreads
	r.values = values
	return
}

func main() {
	world := NewWorld()
	lines := strings.Split(tiny, "\n")
	commands := []Command{}
loop:
	for _, l := range lines {
		c := Command{}
		arr := strings.Split(l, " ")
		switch arr[0] {
		case "on":
			c.action = ON
		case "off":
			c.action = OFF
		default:
			logFatal("Unknown action %q\n", arr[0])
		}
		arr = strings.Split(arr[1], ",")
		for _, data := range arr {
			equals := strings.Split(data, "=")
			numbers := strings.Split(equals[1], "..")
			low, high := strToInt(numbers[0]), strToInt(numbers[1])
			if high < -50 || low > 50 {
				continue loop
			}
			s := Spread{low, high}
			switch equals[0] {
			case "x":
				c.x = s
			case "y":
				c.y = s
			case "z":
				c.z = s
			default:
				logFatal("Unknown param %q\n", equals[0])
			}
		}
		commands = append(commands, c)
	}

	// fmt.Printf("Commands: %+v\n", commands)
	for _, c := range commands {
		fmt.Printf("Processing command %+v\n", c)
		for x := c.x.low; x <= c.x.high; x++ {
			for y := c.y.low; y <= c.y.high; y++ {
				world.Set(Coordinate{x, y, 0}, Spread{c.z.low, c.z.high}, c.action == ON)
			}
		}
		fmt.Printf("Number on: %d\n", world.Len())
	}

	fmt.Printf("Number on: %d\n", world.Len())
}

type Coordinate struct {
	x, y, z int
}

type World struct {
	data map[int]map[int]*RangeMap
}

func NewWorld() World {
	return World{
		data: make(map[int]map[int]*RangeMap),
	}
}

func (w *World) Set(c Coordinate, s Spread, v bool) {
	if _, ok := w.data[c.x]; !ok {
		w.data[c.x] = make(map[int]*RangeMap)
	}
	if _, ok := w.data[c.x][c.y]; !ok {
		w.data[c.x][c.y] = NewRangeMap()
	}
	w.data[c.x][c.y].Add(s, v)
}

func (w *World) Len() uint64 {
	var count uint64
	for _, row := range w.data {
		for _, rm := range row {
			for s, v := range rm.ranges {
				if rm.values[s] == false {
					continue
				}
				dist := v.high - v.low + 1
				count += uint64(dist)
			}
		}
	}
	return count
}

func logFatal(s string, args ...interface{}) {
	fmt.Printf(s, args...)
	os.Exit(1)
}

func minInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func maxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func absInt(a int) int {
	return int(math.Abs(float64(a)))
}

func strToInt(str string) int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logFatal("could not parse %q: %v\n", str, err)
	}
	return int(n)
}

var tiny = `on x=10..12,y=10..12,z=10..12
on x=11..13,y=11..13,z=11..13
off x=9..11,y=9..11,z=9..11
on x=10..10,y=10..10,z=10..10`

var sample = `on x=-20..26,y=-36..17,z=-47..7
on x=-20..33,y=-21..23,z=-26..28
on x=-22..28,y=-29..23,z=-38..16
on x=-46..7,y=-6..46,z=-50..-1
on x=-49..1,y=-3..46,z=-24..28
on x=2..47,y=-22..22,z=-23..27
on x=-27..23,y=-28..26,z=-21..29
on x=-39..5,y=-6..47,z=-3..44
on x=-30..21,y=-8..43,z=-13..34
on x=-22..26,y=-27..20,z=-29..19
off x=-48..-32,y=26..41,z=-47..-37
on x=-12..35,y=6..50,z=-50..-2
off x=-48..-32,y=-32..-16,z=-15..-5
on x=-18..26,y=-33..15,z=-7..46
off x=-40..-22,y=-38..-28,z=23..41
on x=-16..35,y=-41..10,z=-47..6
off x=-32..-23,y=11..30,z=-14..3
on x=-49..-5,y=-3..45,z=-29..18
off x=18..30,y=-20..-8,z=-3..13
on x=-41..9,y=-7..43,z=-33..15
on x=-54112..-39298,y=-85059..-49293,z=-27449..7877
on x=967..23432,y=45373..81175,z=27513..53682`

var input = `on x=-38..7,y=-5..47,z=-4..41
on x=-16..35,y=-21..25,z=-48..5
on x=-43..4,y=-32..21,z=-18..27
on x=-16..38,y=-37..9,z=-10..40
on x=-3..43,y=-40..13,z=-48..-4
on x=-6..43,y=-4..41,z=-6..47
on x=-29..15,y=-9..43,z=-39..5
on x=-37..9,y=-16..37,z=-1..45
on x=-28..21,y=-7..46,z=-10..36
on x=-26..27,y=-6..40,z=-18..34
off x=13..30,y=32..41,z=-10..1
on x=-43..6,y=-7..46,z=-15..31
off x=19..30,y=-43..-27,z=-36..-26
on x=-15..34,y=-41..10,z=-45..0
off x=15..31,y=27..36,z=20..33
on x=-8..42,y=-44..6,z=-22..25
off x=-37..-20,y=22..40,z=35..44
on x=2..46,y=-43..3,z=-17..36
off x=23..34,y=5..16,z=-5..6
on x=-9..36,y=-47..7,z=-47..5
on x=-160..27861,y=57453..76567,z=10007..35491
on x=-34240..-12876,y=-87116..-51990,z=33628..53802
on x=-23262..4944,y=62820..88113,z=-7149..-1858
on x=-48127..-16845,y=59088..86535,z=-5102..744
on x=-54674..-32293,y=21567..47184,z=55761..72099
on x=-69807..-50439,y=-53937..-38431,z=14338..28485
on x=13123..34571,y=-74162..-53699,z=5810..32903
on x=-36751..-16860,y=49545..75658,z=-34660..-21694
on x=-68435..-62588,y=-40370..-6437,z=-39677..-23920
on x=-82974..-59309,y=18720..41775,z=-39850..-17189
on x=-72347..-41708,y=9850..18229,z=-53512..-34975
on x=-30894..4418,y=65635..84979,z=19530..39249
on x=15288..25235,y=-29448..1499,z=71472..85280
on x=-54686..-44312,y=36734..56590,z=-60370..-40555
on x=-5266..21740,y=-63065..-47319,z=-53244..-34997
on x=30363..46391,y=65963..73471,z=-16696..-1815
on x=8243..24278,y=-94018..-69081,z=4245..29009
on x=-37020..-23168,y=-72025..-55051,z=-35968..-10146
on x=-47068..-26645,y=-20779..13147,z=56626..75638
on x=-38276..-11929,y=-75280..-63708,z=19483..24673
on x=3913..16625,y=66848..92341,z=9711..32412
on x=-33513..-7553,y=-44460..-19203,z=53802..84399
on x=-46017..-22850,y=-46589..-32578,z=-70024..-47623
on x=48573..58330,y=2974..29593,z=34454..60055
on x=-10056..10090,y=8646..37530,z=-95368..-64816
on x=-40314..-23851,y=49027..75546,z=33133..59153
on x=9529..39846,y=-43216..-6720,z=66362..84337
on x=-45474..-33807,y=15249..40850,z=-72058..-42686
on x=51316..72889,y=-20375..8377,z=32410..56709
on x=-2105..32624,y=-78827..-73957,z=-36836..-9360
on x=28703..45480,y=7042..25778,z=-82005..-55614
on x=6954..42651,y=4674..37873,z=57207..93287
on x=29735..50615,y=17210..29568,z=64317..74832
on x=-71817..-57182,y=-52937..-34602,z=-16791..8172
on x=62789..73440,y=-8262..6641,z=-47239..-39833
on x=22253..33377,y=-66419..-35421,z=-76817..-56664
on x=-81664..-74289,y=-19781..2191,z=-47797..-9331
on x=-15501..3487,y=57147..69359,z=30481..53449
on x=65898..70034,y=-42224..-28303,z=-36101..-5313
on x=43494..54307,y=42390..54768,z=-48993..-31471
on x=-78098..-51348,y=-31792..-29005,z=-59011..-31680
on x=1733..29690,y=-76915..-65890,z=-51511..-21684
on x=15663..36026,y=-43573..-32068,z=54160..67763
on x=-15046..8020,y=73251..81288,z=-6407..6193
on x=-27828..-7924,y=-81275..-66773,z=-6446..19063
on x=53718..81429,y=8690..21113,z=-46849..-30403
on x=-36511..1151,y=59425..77413,z=7092..45917
on x=1629..15400,y=11693..34093,z=68313..91797
on x=-26483..210,y=1338..31754,z=-84816..-58477
on x=-83237..-62066,y=4539..25394,z=-51906..-28397
on x=-9149..923,y=-81989..-60157,z=-26789..-5365
on x=17517..37608,y=26598..46873,z=57508..64782
on x=-7119..11105,y=62942..83529,z=9625..26513
on x=-51514..-30945,y=20145..42149,z=-62845..-43704
on x=5756..29641,y=-22431..5811,z=63100..78564
on x=2390..19403,y=-20570..2836,z=62003..85329
on x=-30076..-15623,y=55291..75814,z=15603..36585
on x=31384..53207,y=-52113..-34713,z=-53230..-48671
on x=-22679..10340,y=13595..27259,z=74909..82901
on x=67880..82348,y=8971..27422,z=-2124..15487
on x=-18895..10483,y=-85666..-59455,z=34047..51125
on x=55861..76270,y=-63315..-43630,z=-14734..14650
on x=574..29499,y=-45563..-10582,z=-90533..-55006
on x=-69275..-42476,y=-59658..-25077,z=-52699..-28481
on x=-56824..-33936,y=-25868..8029,z=-82735..-53242
on x=-90068..-61105,y=25564..34118,z=-8753..8383
on x=-16748..-8507,y=-60001..-44220,z=39386..65527
on x=47454..61550,y=30183..42751,z=-48341..-40487
on x=-75242..-56708,y=-18345..7252,z=27135..42009
on x=-58608..-41635,y=-77670..-51210,z=-19528..-1735
on x=-52353..-35577,y=65342..71094,z=-34292..-10437
on x=-13904..390,y=-85050..-62408,z=-40671..-16675
on x=-70826..-61780,y=-28261..-8503,z=-43838..-22314
on x=23456..34852,y=-75129..-39146,z=-59367..-27639
on x=-21289..-10174,y=-4853..32679,z=-90593..-71383
on x=-37590..-9368,y=-75394..-60296,z=-44762..-24730
on x=-60024..-32249,y=59509..67462,z=-2182..22377
on x=-8113..11767,y=-35418..-18956,z=64722..86638
on x=46854..82238,y=40215..59083,z=-11037..11812
on x=49339..85055,y=28346..47525,z=5491..39152
on x=-50591..-22806,y=-33668..-17243,z=58370..72406
on x=54558..69946,y=-56413..-35097,z=-27088..7730
on x=-64262..-47635,y=-74110..-35937,z=13887..37783
on x=-51771..-28374,y=-60611..-46568,z=23665..55950
on x=-81107..-54757,y=-26037..7731,z=-40004..-22397
on x=-41417..-19367,y=54161..90953,z=12299..44226
on x=-89805..-55266,y=-31643..-22791,z=-15805..-3577
on x=-78290..-54183,y=8949..18096,z=-47920..-37582
on x=-44972..-27636,y=-56757..-33114,z=43862..61796
on x=34720..58327,y=40641..60954,z=5688..14788
on x=-56983..-26271,y=-69520..-43968,z=5457..30306
on x=52234..80932,y=-28451..-12788,z=25327..31570
on x=-20283..-10624,y=55145..77699,z=-55377..-51731
on x=-8918..23862,y=-26789..-7902,z=78219..83878
on x=-81127..-49088,y=-33284..-13921,z=-50207..-32847
on x=3968..12957,y=62162..90459,z=-32604..-13525
on x=50536..63270,y=47029..71173,z=25852..40402
on x=-19947..3085,y=69295..98557,z=-7457..15697
on x=56416..88951,y=-49125..-23678,z=13033..31254
on x=-45415..-13431,y=-70617..-33184,z=-67667..-43697
on x=-25226..-16517,y=11281..39430,z=65912..87788
on x=32863..54130,y=33778..71445,z=42050..55059
on x=-61047..-29947,y=59749..68086,z=-42477..-13948
on x=-18251..9292,y=-89783..-58004,z=-41787..-15274
on x=5759..16739,y=31155..46337,z=-75036..-53352
on x=31554..54658,y=53243..73463,z=-26575..5972
on x=-33007..-6786,y=-62474..-41434,z=50315..75779
on x=27458..32973,y=22810..54163,z=-75563..-58647
on x=8075..18075,y=-50635..-24172,z=62978..85992
on x=-76736..-47443,y=-4014..19377,z=45970..58134
on x=-60157..-52753,y=-43338..-22023,z=-51719..-41956
on x=-3216..22631,y=49983..61450,z=-77551..-54695
on x=-46004..-40643,y=51830..62022,z=18230..48510
on x=42330..66450,y=-40592..-14557,z=-58845..-29636
on x=-32729..-9787,y=-91565..-69806,z=18117..35690
on x=-32002..-10885,y=-1443..14786,z=-90524..-68253
on x=50978..75653,y=7140..18805,z=-63220..-43305
on x=-64568..-49330,y=15278..36078,z=28629..42531
on x=19040..44557,y=-75487..-59006,z=-37514..-20407
on x=-2797..6071,y=-87782..-64015,z=-23722..-9160
on x=30427..40250,y=-51346..-20826,z=-76885..-46217
on x=-9837..5654,y=49778..66163,z=42263..63916
on x=-21072..11872,y=28390..51048,z=-77577..-49183
on x=31541..48623,y=37624..58494,z=-42361..-30931
on x=-56299..-27988,y=47361..70101,z=1213..10546
on x=-23456..-16890,y=-79546..-59160,z=-6743..10289
on x=14204..41142,y=56511..74601,z=5095..33434
on x=-36005..-14701,y=-68016..-39674,z=36553..63480
on x=60545..74232,y=15015..42503,z=-43149..-34200
on x=41594..65337,y=50002..72588,z=12856..37989
on x=-90432..-62701,y=-37837..-4386,z=-32554..-13392
on x=70315..83165,y=-40717..-23351,z=-28440..-1511
on x=-24120..5947,y=34921..55704,z=-69582..-46555
on x=-32209..-8756,y=-77077..-58552,z=-48224..-14837
on x=-43510..-17400,y=64412..85528,z=6143..27595
on x=36053..55013,y=53314..68720,z=-22332..-5189
on x=-29743..-3821,y=65395..78815,z=-9292..23641
on x=-69714..-38553,y=-44109..-30469,z=29793..51760
on x=62579..74277,y=11751..29748,z=23050..34754
on x=-2338..12874,y=75847..79442,z=10731..41551
on x=-47619..-35145,y=-52665..-37628,z=39702..49423
on x=-18665..-6319,y=-71514..-63562,z=42460..45818
on x=-30044..-12170,y=53793..59153,z=-63499..-40230
on x=-12567..16134,y=-90016..-66446,z=-35676..-19472
on x=-85062..-55201,y=4912..28423,z=-41593..-11980
on x=56423..80218,y=-37562..-16517,z=36700..48100
on x=-79818..-55704,y=11491..34326,z=8367..40952
on x=-26293..-12220,y=41552..57647,z=51547..73246
on x=-92036..-71546,y=10557..30361,z=-37285..-11767
on x=36855..65165,y=-77057..-59070,z=4400..24300
on x=-50136..-35853,y=35165..54382,z=37467..66189
on x=16047..55776,y=-39553..-11227,z=58349..80407
on x=27458..35035,y=-32489..-4540,z=-78886..-68473
on x=46937..76452,y=12250..27089,z=-55190..-45467
on x=-60706..-29052,y=-43190..-26027,z=48544..56051
on x=33102..61073,y=27580..41152,z=41873..64266
on x=-83817..-62138,y=-874..20403,z=-57292..-40933
on x=54691..83982,y=-26668..-1254,z=25108..50265
on x=-64405..-32749,y=44513..79090,z=-13114..-3445
on x=4627..28855,y=-38014..-18427,z=-78981..-56867
on x=39450..53446,y=42289..75373,z=-33901..-6131
on x=63899..90602,y=26653..55063,z=-10988..9550
on x=14634..44509,y=-78079..-56702,z=28926..52784
on x=38612..58844,y=41733..50901,z=-69067..-31286
on x=-76929..-61757,y=-65663..-31719,z=-735..7695
on x=28516..47164,y=-29956..-15347,z=-75568..-61002
on x=-50167..-42193,y=-40926..-30244,z=-60920..-47646
on x=-55502..-47750,y=54875..64454,z=-29489..-3798
on x=14430..36601,y=26459..39911,z=-79289..-57027
on x=25075..27283,y=52962..87505,z=18854..33830
on x=-58860..-43502,y=58845..85498,z=-18144..7960
on x=-53792..-49597,y=-62122..-48441,z=-6272..16303
on x=-1823..9424,y=-44705..-17073,z=70120..80148
on x=63692..82216,y=-19361..8359,z=33210..45303
on x=-11840..-2979,y=24389..52039,z=60315..71939
on x=-73992..-63729,y=-45246..-14461,z=-46112..-19286
on x=44923..62221,y=-54960..-42011,z=-38137..-20135
on x=14722..33976,y=-61145..-30720,z=-75360..-54000
on x=245..27732,y=58258..82234,z=30097..48740
on x=6991..26059,y=17174..25794,z=71284..89533
on x=35195..56084,y=51212..74918,z=-32483..-9615
on x=33703..66981,y=-60122..-54236,z=-53453..-17926
on x=-61940..-36371,y=-16581..-4570,z=-60706..-40703
on x=60760..77773,y=-33429..-9006,z=14598..31427
on x=18431..39521,y=-7515..9629,z=60150..89625
on x=54831..71507,y=20590..53690,z=18166..24092
on x=-25927..-4852,y=56797..80031,z=24590..44494
on x=58156..74392,y=-54667..-25937,z=-34893..-13568
on x=-46251..-30782,y=70616..92727,z=-10301..13242
on x=-54435..-31583,y=18809..57238,z=52002..64018
on x=-87191..-65448,y=21016..30554,z=-7539..18577
on x=-27162..-8066,y=-76175..-58675,z=-45390..-18705
on x=-28095..-21402,y=-70382..-55725,z=38509..56620
on x=49521..65654,y=31848..57013,z=29190..51381
on x=-75558..-48003,y=15704..38247,z=17226..43613
on x=23884..58026,y=50160..65808,z=6785..27572
on x=58063..93764,y=-22911..-4938,z=-11895..2739
on x=43350..68759,y=-71427..-46464,z=-6078..15892
on x=-33140..-3203,y=-20673..-10382,z=-75944..-67085
on x=-21259..-17037,y=67489..78783,z=-8541..11664
on x=36779..59339,y=24261..52940,z=-53125..-35093
on x=-30493..793,y=44829..70216,z=-52950..-47982
on x=45292..72628,y=-55273..-48706,z=-1598..9919
on x=29403..34110,y=-54445..-36095,z=43199..56783
on x=-59861..-33784,y=-53249..-37265,z=34591..65790
off x=41754..57148,y=-80200..-55875,z=3599..22722
on x=-10075..4386,y=27944..45112,z=-79226..-60131
off x=-15187..248,y=-83279..-72364,z=5037..32511
on x=54933..58881,y=12281..29091,z=-69028..-54121
on x=8090..25240,y=70064..81196,z=-16223..-13021
off x=-60400..-26105,y=-62607..-40191,z=47851..67422
off x=24451..46164,y=-75961..-52303,z=-49916..-25476
on x=-44857..-7989,y=-38916..-18382,z=-70680..-58913
off x=-72200..-43203,y=24727..57702,z=-44696..-7928
on x=26323..57990,y=-14888..994,z=-79264..-65525
on x=-48909..-40386,y=-73244..-64696,z=-323..16342
on x=35358..52991,y=10687..27904,z=-70944..-57912
off x=-23236..-4162,y=-29712..-17951,z=59106..79389
off x=53581..72320,y=6800..13298,z=-45872..-42480
on x=-43571..-35441,y=41386..68207,z=-44478..-35187
on x=29214..36514,y=2338..22755,z=54949..83144
on x=-84002..-51122,y=41833..51908,z=9476..15499
off x=7540..23221,y=-12461..-7794,z=-95431..-66096
on x=76564..93348,y=-9452..16724,z=-13630..7535
on x=24899..35662,y=-61799..-36371,z=50395..73121
off x=-23379..-14671,y=-61654..-26944,z=-66602..-52265
on x=-75965..-53958,y=-25180..-13183,z=24409..52392
on x=-10197..4622,y=-15474..1783,z=-92565..-76011
off x=45082..71453,y=9055..32401,z=53858..58811
on x=-76326..-55011,y=-39530..-20949,z=21859..44714
off x=32466..57949,y=29661..50902,z=38020..60279
off x=-48739..-19365,y=36215..74060,z=-61189..-33437
off x=-13146..-2771,y=9021..18827,z=60320..97125
on x=-224..15138,y=31517..57835,z=-83757..-53030
on x=48425..68878,y=14971..40195,z=-64239..-33895
on x=-3608..5152,y=-30571..-5399,z=75733..96504
on x=62972..67090,y=-52211..-34797,z=-11592..8052
on x=6468..27085,y=-75658..-53084,z=24332..51031
on x=-9175..29042,y=-47119..-33884,z=-82475..-62552
off x=-15847..4062,y=60344..89283,z=-30211..-17239
on x=25768..49899,y=62273..71708,z=349..10734
off x=42707..63210,y=13025..46778,z=36022..55110
on x=46293..67516,y=-7202..26002,z=-66272..-31228
off x=-62283..-34594,y=-38660..-16958,z=38044..60467
off x=-70095..-49102,y=31813..56709,z=8713..28863
on x=13788..17456,y=-12132..6791,z=61871..77520
on x=-43439..-29707,y=-56265..-39662,z=-72104..-38103
off x=13592..44088,y=-74894..-45793,z=37968..66601
off x=46700..52734,y=-40463..-19848,z=-69191..-42886
on x=20538..39375,y=61182..85768,z=15326..29409
on x=54417..76403,y=-19878..-1308,z=-42875..-26138
on x=-86317..-61464,y=-18403..3486,z=-18507..6337
on x=2786..19873,y=-20405..9939,z=59315..81073
off x=-59877..-54443,y=34487..69275,z=3488..18241
off x=-39007..-20035,y=4078..30312,z=-82731..-51587
on x=32588..56936,y=-76500..-67396,z=-21385..-397
off x=25080..41990,y=-28846..-19737,z=69970..78086
on x=23152..34879,y=52174..69560,z=27953..49114
on x=60684..71642,y=-30692..-3159,z=31270..54443
off x=-45205..-29098,y=58510..89307,z=-17925..2596
off x=19582..24130,y=-83134..-57590,z=35254..46862
on x=-11911..23624,y=-67005..-30930,z=-81907..-60918
off x=4727..24751,y=-39077..-37491,z=-71684..-61996
off x=-9593..7610,y=-9049..14670,z=62634..98129
on x=-66988..-43094,y=-62996..-35119,z=11118..43202
on x=-35998..-5275,y=-37152..-15041,z=-87055..-60774
off x=-14073..-1453,y=48138..72420,z=-55577..-47324
on x=53303..57735,y=-50475..-39459,z=28711..49776
on x=50327..79524,y=-20647..-7189,z=32930..64357
off x=-21585..-2673,y=67475..97316,z=-16225..10334
on x=12761..42162,y=-75198..-45687,z=27773..56590
off x=10780..29984,y=30806..53252,z=45151..82121
on x=46365..68175,y=-3920..2578,z=-68115..-55420
off x=26455..49250,y=-55979..-42110,z=-63322..-50217
off x=-62311..-44896,y=-70433..-49939,z=-6372..6167
off x=-68448..-49911,y=12490..21799,z=-66820..-50662
on x=11680..41135,y=42991..76600,z=38991..56525
off x=60371..92603,y=19664..34168,z=-32784..-13415
off x=-72850..-41597,y=-58644..-37918,z=-37165..-19623
off x=53858..74510,y=34819..52123,z=16497..43908
on x=48488..52121,y=-45360..-24737,z=-66885..-37749
on x=20448..57476,y=27537..48980,z=46019..57068
on x=-27248..-19724,y=41351..47693,z=-77216..-44436
on x=-52950..-33175,y=41728..62787,z=-16516..-11981
off x=-51177..-31463,y=-63233..-45533,z=32418..61549
off x=-37057..-12370,y=18117..43449,z=-74085..-60338
off x=-12503..2846,y=-42876..-30104,z=63819..75029
on x=60667..70952,y=-46497..-30751,z=16478..34217
on x=10069..26074,y=50365..61240,z=-68598..-56464
off x=17685..45785,y=44978..64857,z=53362..71129
off x=-73990..-42336,y=-72284..-40294,z=8863..26728
off x=68782..80879,y=22420..46148,z=-2329..26525
on x=59344..83826,y=-40162..-30633,z=2097..16442
on x=-27562..-10414,y=27820..58187,z=-65865..-46977
on x=36887..65981,y=-73709..-52311,z=-42786..-24070
off x=-20682..-3178,y=64806..78374,z=10659..34470
off x=-8804..15550,y=-8444..-2404,z=-96231..-65363
off x=-66648..-53156,y=-56554..-29056,z=-33726..-17793
on x=-203..26059,y=68129..80958,z=-43219..-27404
on x=47218..72444,y=35342..44493,z=-36229..-12760
on x=46904..60309,y=-54993..-32818,z=-50008..-30454
off x=-21694..-6209,y=-88469..-65433,z=-11873..12245
off x=-62399..-24805,y=28232..34853,z=51230..60419
on x=-79596..-57362,y=-48983..-34721,z=-13420..10701
on x=-6768..1790,y=-19645..-5406,z=-92627..-74287
on x=-33652..-18284,y=-23439..-8382,z=67655..90006
off x=16524..39458,y=29152..49162,z=55749..83384
on x=-54254..-34002,y=-88855..-70903,z=-22296..3410
on x=-20779..8521,y=-89181..-73450,z=-7812..26787
on x=1615..21625,y=-6080..17570,z=71668..83316
off x=-9715..12709,y=-50780..-22669,z=-74774..-69626
on x=9867..24037,y=12358..42384,z=72370..82642
on x=3224..30208,y=28864..51571,z=61883..84706
off x=-945..19284,y=66805..89823,z=25957..49425
off x=-25280..-6091,y=-97225..-62630,z=-19807..-7788
off x=-2250..13009,y=50093..83251,z=-51760..-30168
off x=-54440..-34334,y=26817..48120,z=48074..60160
on x=23262..49261,y=40455..52393,z=-60959..-43615
on x=-23294..-2268,y=-52310..-24691,z=-87802..-57989
off x=-26612..7840,y=72808..79999,z=-29713..-3816
on x=54391..88234,y=-51754..-15468,z=1945..19324
off x=17508..32841,y=-73524..-65420,z=-35218..-11372
off x=31152..49478,y=-70736..-54693,z=-43913..-40242
on x=53480..64046,y=-62586..-42748,z=34500..47779
off x=34058..63228,y=-61928..-58936,z=-31573..-3659
off x=-58697..-39679,y=-81148..-49181,z=13624..23687
off x=-27245..-8270,y=35092..58039,z=40192..57884
off x=32739..63309,y=-40682..-17682,z=47869..70035
on x=52219..64202,y=-67303..-51407,z=4800..23668
on x=-54654..-29433,y=-86817..-62196,z=-5556..23658
off x=-75279..-49504,y=-7178..9542,z=44327..53956
on x=-7810..8890,y=-30024..-89,z=66319..97476
on x=8400..36435,y=12009..44304,z=-80097..-67512
off x=22016..38709,y=71634..74972,z=-3129..1171
off x=48394..68796,y=44368..57057,z=-40373..-16910
off x=-56815..-38268,y=36694..67658,z=-47750..-32089
off x=-57877..-35052,y=46097..56601,z=33355..49073
on x=23011..37434,y=-53974..-23058,z=-61883..-45402
off x=11872..42043,y=19031..22704,z=-83665..-55245
off x=-73243..-46798,y=-4329..17893,z=-57779..-44793
on x=67456..73394,y=-22388..-14950,z=-37499..-28618
off x=30365..66533,y=-70682..-46579,z=-16208..-5031
on x=21360..45088,y=-76539..-50147,z=21845..41905
on x=-18419..-11559,y=76149..84053,z=-11328..-2185
off x=33020..55476,y=-15536..2658,z=64275..79458
on x=-48063..-26857,y=-20734..10735,z=-85808..-52673
off x=43790..66834,y=-56342..-29394,z=31254..58327
off x=67314..83024,y=-13230..17187,z=-23371..-4720
off x=-40007..-14810,y=-73067..-57103,z=-66195..-38082
off x=25841..57188,y=-26024..897,z=-87581..-52851
off x=45794..76267,y=-52081..-18959,z=-40107..-15901
off x=-40509..-4365,y=56716..81562,z=-22853..6830
off x=-77503..-45958,y=-63620..-34287,z=-19880..-7413
off x=-19264..12115,y=-19022..8117,z=-93523..-72932
off x=39766..60651,y=-59973..-42774,z=-28355..-7554
on x=-5185..26343,y=42329..68147,z=56719..69801
off x=-1117..14361,y=41023..50863,z=-73482..-62392
off x=14618..41995,y=-57249..-23838,z=-82100..-56433
on x=-73588..-51361,y=3468..30508,z=-52360..-20267
on x=12744..37529,y=-54963..-50978,z=-53987..-45070
off x=-72308..-43455,y=-38389..-24377,z=-50714..-45259
off x=-78117..-42852,y=-13721..13148,z=-55905..-37080
off x=-4467..22317,y=-12170..5957,z=71290..88500
on x=-95981..-66408,y=7590..8037,z=-2404..26148
off x=12232..45401,y=12849..45079,z=-78972..-62278
on x=-46581..-20385,y=43824..70561,z=32388..58000
off x=14142..24293,y=-26356..7327,z=-82833..-56772
on x=-62185..-54635,y=-62453..-44051,z=-17689..-15452
off x=7493..30003,y=51943..75363,z=19045..35329
on x=10852..22754,y=-21936..-12574,z=58047..81390
on x=-85885..-60870,y=-5186..26456,z=-18838..-15144
off x=573..10699,y=-12795..12158,z=62024..82285
on x=-53647..-41451,y=6744..19177,z=-65585..-52437
on x=-1937..13101,y=16787..44245,z=-89301..-69691
off x=11648..28117,y=65897..81495,z=-51744..-31456
on x=-34510..-23114,y=64316..84122,z=11745..26248
on x=46597..67036,y=-49078..-18339,z=-55106..-43379
on x=-1630..16917,y=44609..68512,z=50592..75476
off x=-21971..1123,y=24005..48686,z=60732..80668
off x=42871..68726,y=-60814..-42185,z=27184..53921
off x=38826..47043,y=-40691..-34285,z=-57815..-35188
on x=44616..56936,y=-18905..4924,z=53303..63015
on x=22934..47063,y=21220..45110,z=48969..72456
off x=31507..57059,y=41717..67987,z=-48794..-42062
on x=63949..80604,y=11054..40766,z=17043..53534
on x=28112..49371,y=-69150..-55767,z=-37760..-7539
off x=-39198..-18198,y=-10040..11413,z=-84459..-58961
on x=-21052..258,y=-11980..13528,z=-92447..-67079
off x=-72190..-65264,y=-11322..11452,z=20127..52589
on x=9456..34476,y=75111..90133,z=-3087..22015
off x=-48678..-16972,y=-41882..-11490,z=63832..72591
on x=711..33910,y=-92450..-65393,z=-19051..4525
off x=62421..84065,y=26556..40786,z=14281..29967
on x=-67588..-56095,y=28175..65003,z=22332..47033
on x=-70974..-36736,y=9539..20776,z=49351..68205
off x=-25177..13122,y=-61323..-41425,z=46642..75842
on x=-52446..-40486,y=-62884..-38244,z=20759..52144
on x=-42548..-33936,y=-53496..-46698,z=-50788..-30250
off x=33339..50697,y=-75051..-48370,z=36655..54219
off x=43213..72026,y=-22655..-11856,z=37238..51298`
