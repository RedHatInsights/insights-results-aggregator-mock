/*
Copyright © 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package data

import (
	"github.com/RedHatInsights/insights-results-aggregator-mock/types"
)

// RequestIDs contains mapping between cluster name and sequence of request IDs
var RequestIDs map[types.ClusterName][]types.RequestID

// init is called before the program enters the main function, so it is perfect
// time to initialize maps etc.
func init() {
	RequestIDs = make(map[types.ClusterName][]types.RequestID, 10)

	// up to 12 request IDs per cluster
	RequestIDs["34c3ecc5-624a-49a5-bab8-4fdc5e51a266"] = []types.RequestID{
		"3nl2vda87ld6e3s25jlk7n2dna",
		"18njbjudvkc521w8buicx2clri",
		"38584huk209q82uhl8md5gsdxr",
		"1wpkdrz1gknwe3p2ymia8g4s97",
		"17xqaht4cxigc1y0w87k47ygbh",
		"21hog6cmurvyw3m7kyqczt3hgs",
		"39p67r9p5hqz43f8k7kugm1zgn",
		"1cwa2ehvh9dnl2pyuh6kfrdrjl",
		"3gas724nbyy3l1rx77h6ewsi0y",
		"16z7mxc44n1ol1cr5um7vxgpqh",
		"1ccu11bn9ek2x20ink95zo0efg",
		"3oeiljuhkvbi61hf6tpgk4p2sk",
	}
	RequestIDs["34c3ecc5-624a-49a5-bab8-4fdc5e51a267"] = []types.RequestID{
		"1duzaoao0l1b230ipv0rb4sqe8",
		"1yjdje758zgyy3ksfr732yb1cl",
		"2ukce6u74rm9e12mplu6liqzsv",
		"2mie19y6yt6xo1lhyrfw7pea6h",
		"1xvgh46mtvj9b2g6qpwvkqf4ay",
		"2gb0ayi7x579e2noepkpbqyoa8",
		"2i66cwa3ee0zu39z8aeg3ucu81",
		"2dd2h2rbt8mi03gn0jugwp8ybs",
		"2b0vbkdy2t3u12susperilgmee",
		"2drtvjlisiqww1c93kugqyboyc",
	}
	RequestIDs["34c3ecc5-624a-49a5-bab8-4fdc5e51a268"] = []types.RequestID{
		"3ga8nnps8unei12m20neiqchmd",
		"2nmk76dfs1n7h1ua47bb7f21ht",
		"2anbyo69luzzq3vltwpqj4dbog",
		"2f3cvuy851bd21rnt0t4tzpv7r",
		"2q8rrjoyufzl23mfc6x2p6e875",
		"3fninx6jrt9un2wghc3rd59w0d",
		"16hhloic626c41uunfyy62wc34",
	}
	RequestIDs["34c3ecc5-624a-49a5-bab8-4fdc5e51a269"] = []types.RequestID{
		"1l2z27ytw9pbd2tsjz3to6amu0",
		"2im6w4fhn0els11wz5yuvsjr8q",
		"2gixiqa8wmxod2h5i5dot9dcrw",
		"32hbh1itylw9l37b145e180gvf",
		"1fsexserm84ew2aievqmf89cj0",
		"2kx5vv6cdpgl42qioi8ya91auj",
		"2pgixxxlqksrl3ckgeza0smx4p",
		"204x056xc2cog2hufbnmhevg9v",
		"2rfj50xipfj431ynixwl08z1ky",
		"2y01d1fz9mymf2aj0rfhjbjsp8",
		"2m8fw5lpnkhjc13fv5l5d0ziwt",
	}
	RequestIDs["eeeeeeee-eeee-eeee-eeee-000000000001"] = []types.RequestID{}
	RequestIDs["00000001-eeee-eeee-eeee-000000000001"] = []types.RequestID{}
	RequestIDs["00000003-eeee-eeee-eeee-000000000001"] = []types.RequestID{}
	RequestIDs["74ae54aa-6577-4e80-85e7-697cb646ff37"] = []types.RequestID{
		"1zlcewj4kjtsp37x0yyr6cwhgr",
		"3m3imli92shw225d4c3glzycxq",
		"13yqlst6dmdji2z717w2v5fwcp",
		"271w1b53jlfjq2axaetgpe0yrd",
		"32zr43d2a4cbq1ogi1eu3hrti1",
		"3pyjpvp4umqwx1xnhdq3mwgzkh",
		"33mk6kyt1sml41dzi8ey2rp6s3",
		"3aiwib6toy56o1lpq9g1paciq8",
		"173zl1rkkbnwp2loeest0c9lju",
		"1vb0ozrsj5i8c2v28ajrw0l0gp",
		"1tm87s8lgxk6s2vuixu6gp48t3",
		"2msu764i6zctg1ujo4n5j8jixo",
		"1sfxeleizlijx3nl72qdojlqtm",
		"2bbe40wrkf4fg1stf5cyb3s7ol",
		"24hjylbhu60um37eevp98099ib",
		"2vonepjvxwtuo212z88m9kf7e2",
		"1h3y9oexbese916l327k83rzxk",
		"1trr8kq9audlu1xx59q3l2exw2",
		"2a0tiy6bhgjlm1rx4do0wu018w",
		"2yrlu653bbcp21perspbd20bj1",
	}
	RequestIDs["a7467445-8d6a-43cc-b82c-7007664bdf69"] = []types.RequestID{
		"1bv60sffeojsb3065r2zip4rxu",
		"350nfzi2lwyw11hx8aihf55xnl",
		"2l44bdrtug3x72m5cz1yf3yl52",
		"1k02qxcy6e3y515aowk1vmyolx",
		"2tj3ibnydkm362ssxohrncd82n",
		"16ox0nc4fwmzq2k6cylncr0f16",
		"25iaksaj2n3kj2z48xqjn16uhs",
		"1pmisu4bznd0x2stkh7a7qp2s1",
		"2gxbldelvzque14au4bspza407",
		"37hv36hns7nxl2pwoybnxohz3c",
		"35gxiqmw00wda39bdg8npncy6i",
		"1mnug872sbcrt15ajyi977zmwm",
		"1kfr26jjxy4l72rpgrnn4g0egj",
		"21h34bv52kji82vaj7yrxskxnt",
		"3rruw7n0w9mdy2aw1br2erqs4q",
		"2i3acyyq6opc02rivqaavzwacb",
		"2rjzrh18g3p3e2ajflf0cpa7f7",
		"2p4f7om16ez5s3854qtfsndnwg",
		"2htin0agkfq3g1t74b9dw0q42v",
		"21cxbb3w7vngc199djupcenxcn",
	}
	RequestIDs["ee7d2bf4-8933-4a3a-8634-3328fe806e08"] = []types.RequestID{
		"1kg2gt2u2v51q32s6wbmzbmpix",
		"2q0de4r7msgkr2lo2hq6zg793t",
		"29kmeg57zf6381a49w53uh9wsf",
		"1j1flqu0iaak136xa7afw0z0dp",
		"221ecf70cpvao2d00afur6zjvc",
		"31hyt1tzwrvda27tx60hec1wd7",
		"27iyvyi4fcqid2ffbrp08r61pb",
		"3s1wzfbsnwz5t23z17mj4rrxky",
		"267idrw6e2y6n3bkykurewxp6l",
		"1oz8ugidxc4nr3c3ug59ugp4co",
		"3t4cda6pbykko1zua4lbucm20t",
		"2ewgmnuikkaux3b53v6rl3h6fj",
		"2mycihdys646x14drjpb3yofwi",
		"14y1d4n47cohr3nt4i2h561blo",
		"28oh6m8t5gj7v3g718j5b3yzv3",
		"2zlewqf2pa9za2bfw4vhtrzqta",
		"1zrkw4yip6bin15fp2j9vzv2ak",
		"2t44ix0224bif2ggsb2a99qg58",
		"2rpknxj10r4cq127dkni75lrf2",
		"3ac2ctg26hxcs1ta99lxn64qir",
	}
}
