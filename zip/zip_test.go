package zip

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
	"bufio"
)

var testHtml =`
<!DOCTYPE html>
<html>
<head>
<link rel="dns-prefetch" href="http://i.tq121.com.cn">
<meta charset="utf-8" />
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
<title>【北京天气】北京今天天气预报,今天,今天天气,7天,15天天气预报,天气预报一周,天气预报15天查询</title>
<meta http-equiv="Content-Language" content="zh-cn">
<meta name="keywords" content="北京天气预报,北京今日天气,北京周末天气,北京一周天气预报,北京15日天气预报,北京40日天气预报" />
<meta name="description" content="北京天气预报，及时准确发布中央气象台天气信息，便捷查询北京今日天气，北京周末天气，北京一周天气预报，北京15日天气预报，北京40日天气预报，北京天气预报还提供北京各区县的生活指数、健康指数、交通指数、旅游指数，及时发布北京气象预警信号、各类气象资讯。" />
<meta name="msapplication-task" content="name=天气资讯;action-uri=http://www.weather.com.cn/news/index.shtml;icon-uri=http://www.weather.com.cn/favicon.ico" />
<meta name="msapplication-task" content="name=生活天气;action-uri=http://www.weather.com.cn/life/index.shtml;icon-uri=http://www.weather.com.cn/favicon.ico" />
<meta name="msapplication-task" content="name=气象科普;action-uri=http://www.weather.com.cn/science/index.shtml;icon-uri=http://www.weather.com.cn/favicon.ico" />
<meta name="msapplication-task" content="name=灾害预警;action-uri=http://www.weather.com.cn/alarm/index.shtml;icon-uri=http://www.weather.com.cn/favicon.ico" />
<meta name="msapplication-task" content="name=旅游天气;action-uri=http://www.weather.com.cn/trip/index.shtml;icon-uri=http://www.weather.com.cn/favicon.ico" />
<script src="http://dup.baidustatic.com/js/ds.js"></script>
</head>
<body>
<input id="colorid" type="hidden" value="预报">
<script type="text/javascript" src="http://i.tq121.com.cn/j/weather2014/rili.js?id=201511"></script>
<style>
body{background-color:#fff;background-attachment:fixed;background-position:center top;background-repeat:no-repeat;padding:0;margin:0;border:0;font:14px "Microsoft Yahei",Tahoma,SimSun}*{padding:0;margin:0;}input{outline:none;}input::-ms-clear {display: none;width : 0;height: 0;}img{border:0}ul{list-style:none}a{text-decoration:none;color:#252525}a:hover{color:#ee842f;text-decoration:none}em{font-style:normal}.fl{float:left;display:inline}.fr{display:inline;float:right}.clear{clear:both}.line{border-bottom:1px solid #d7d7d7;height:0;overflow:hidden;width:1000px;clear:both;float:left}.bottom-box,.footer,.header-box,.menu-box,.sheng-show,.top-box{width:100%;min-width:1000px;clear:both}.bottom,.dl-box,.header,.main,.menu,.nav,.top{width:1000px;margin:0 auto}.top-box{line-height:25px;height:25px;border-bottom:2px solid #dbe7f3;background:#f5fafe;font-size:12px}.top a,.top span{color:#252525;margin-right:20px}.top a.en{margin-right:0}.header{height:45px}.bottom-top dl dd a:hover,.city a:hover,.menu a:hover,.top a:hover,.warning a:hover{color:#ee842f}.header-box{background:#f5fafe;z-index:1005;padding-top:20px;position:relative}.menu{display:none}.menu a{font-size:15px}.menu a.color{color:#ee842f}.menu-box .sheng{display:block}.search-box{margin:0 50px;width:450px;position:relative;margin-right:0}.search{width:457px;position:relative}.search input#txtZip{font-family: 微软雅黑;margin:0;color:#aaa;text-indent:0;width:330px;padding:0 30px 0 10px;font-size:15px;height:35px;line-height:35px;background:url(http://i.tq121.com.cn/i/weather2015/city/jt-b.png) no-repeat 350px center #fff;border:1px solid #6eafd7;height:35px;border-radius:2px}.search input#txtZip:hover{border:1px solid #f68227;}.search input#btnZip{  font-family: 微软雅黑;text-indent:0;width:80px;height:37px;background:#f68227;color:#fff;margin-left:3px;cursor:pointer;font-size:16px}.search input:hover,.search.hover{border:1px solid #f68227;}.city-box{z-index:13;position:absolute;left:0;top:43px;border:1px solid #abaeaf;width:370px;display:none;background:#fff;box-shadow:4px 1px 9px -3px #888;height:auto}#wrapper{padding-top:10px}.weatherwapper{padding-top:10px}.weatherMain{margin-top:10px}.box{width:1000px;margin:0 auto;}.city-tt{clear:both;position:relative;overflow:hidden;background:#f6fcff;background: #f6fcff url(http://i.tq121.com.cn/i/weather2015/city/bj-dian.png) repeat-x left bottom;}.city-tt a{width:95px;overflow:hidden;display:block;float:left;border-right:1px solid #d4dde5;font-size:15px;text-align:center;background:#f6fcff;color:#252525;line-height:30px;border-bottom: 1px solid #d4dde5;}.weaper{margin-top:10px}.city-tt a:hover{color:#ee842f}.city-tt b{position:absolute;top:10px;display:block;cursor:pointer;right:10px;width:12px;height:12px;background:url(http://i.tq121.com.cn/i/weather2015/city/iconalls.png) no-repeat -64px -646px}.city-tt a.cur{background:#fff;border-bottom:none;padding-bottom:0;color:#ff6400;  border-bottom: 1px solid #fff;}.city-tt span{width:95px;overflow:hidden;display:block;float:left;border-right:1px solid #d4dde5;font-size:15px;text-align:center;background:#f6fcff;color:#252525;line-height:30px}.city-tt span.cur{background:#fff;border-bottom:none;padding-bottom:0;color:#ff6400}.w_city{overflow:hidden;clear:both;padding:10px 0;display:none}.w_city a{width:59px;height:20px;line-height:20px;padding:0;float:left;text-align:left;font-size:14px;margin-bottom:7px;color:#252525}.gn a{padding-right:30px;width:auto;float:left;text-align:left;font-size:14px;margin-bottom:7px;color:#252525}.city_guonei dl dd.diq a{padding:0;width:62px;float:left;text-align:left;margin-bottom:7px}.w_weather{float:right;width:auto;line-height:43px;position:relative;margin-top:5px;cursor: pointer;}.w_weather span{color:#7b7b7b;display:block;float:left;font-size:17px;line-height:30px;height:30px;width:68px;margin-right:15px;overflow:hidden; text-align:right; position:relative;}.w_weather span em{font-size:17px; display:block; white-space:nowrap; position:relative; left:0; top:0px;}.w_weather span em:hover{color:#ee842f}.w_weather em.s{line-height:30px;margin:0 10px;color:#0e70a1;font-size:16px}.w_weather a img{margin-top:14px}.w_weather a.add{background:url(http://i.tq121.com.cn/i/weather2015/city/iconalls.png) no-repeat -11px -593px;color:#0e70a1;font-size:25px;height:32px;margin-top:8px;padding:0px;text-indent:-9999px;width:20px}.w_weather a.dz_right{background-position:-61px -593px}.w_weather a.dz_down{background-position:-35px -593px}.w_weather .more{position:absolute;right:0;top:20px;width:278px;display:none}.nav{margin-top:20px;height:33px;padding-bottom:10px;color:#252525}.nav a{padding:0 52px;font-size:20px;color:#252525}.nav a.sheng{background:url(http://i.tq121.com.cn/i/weather2015/city/iconalls.png) no-repeat 94px -640px;color:#ee842f}.sheng-show{background:#ebebeb;display:none;position:absolute;top:118px;left:0}.dl-box dl{padding:10px 0}.dl-box dl dt{background:#076ea9;float:left;font-size:16px;height:40px;line-height:20px;margin-right:4px;padding:13px 10px;text-align:center;width:46px;border-radius:50%}.dl-box dl dt a{color:#fff}.dl-box dl dd{float:left;font-size:16px;height:56px;line-height:28px;padding:5px 0}.dl-box dl dd a{color:#076ea8;padding:0 16px}.menu-box{background:#e6f4ff;height:40px;line-height:40px;margin-top:1px;border-bottom:3px solid #c3daec;margin-bottom:10px}.menu{text-align:center;font-size:15px}.menu a{margin:0 37px;display:inline-block;color:#252525}.city_guonei dl{width:360px;height:auto;float:left;border-bottom:1px dotted #ccc;margin:0 auto;margin-left:4px;margin-bottom:8px}.city_guonei dl dt{float:left;text-align:center;font-size:14px;color:#252525;width:60px}.city_guonei dl dd{width:299px;float:right}.city_guonei dl dd a{font-size:12px;color:#6f6f6f;word-wrap:break-word}.city_guonei dl dd a:hover{color:#ee842f}.w_logo a{background:url(http://i.tq121.com.cn/i/weather2015/city/iconalls.png) no-repeat 0 -476px;display:block;width:172px;height:43px}.w_yj{margin-right:8px;width:23px;height:18px;background:url(http://i.tq121.com.cn/i/weather2015/city/iconalls.png) no-repeat -177px -501px;margin-top:12px}.w_weather .more{padding:0;z-index:999;border:1px solid #dcdada;position:absolute;right:0;top:40px;background:#fff;box-shadow:4px 1px 9px -3px #888;display:none}.w_weather .more li{height:20px;line-height:20px;position:relative;border-bottom:1px dashed #d6dbe1;width:260px;margin:0 auto;padding:10px 0}.w_weather .more li.on{background:none repeat scroll 0 0 #ebebeb;border:0}.w_weather .more li.on i{background-image:url(http://i.tq121.com.cn/i/weather2014/jpg/blue17d.jpg?v=1)}.w_weather .more li a{color:#669ec0;cursor:pointer;display:block;font-size:13px;height:40px;position:absolute;left:0;top:0;line-height:40px}.w_weather .more li span{line-height:39px; float:left;text-align:left;display:block;overflow:hidden;height:39px;margin-right:6px;font-size:12px;color:#043567;width:70px;text-align:center}.w_weather .more li i{background-image:url(http://i.tq121.com.cn/i/weather2014/jpg/blue17.jpg?v=1);float:left;height:17px;margin:3px 2px;width:17px}.w_weather .more li em{float:left;text-align:right;width:69px;color:#043567;font-size:12px}.w_weather .more li a:hover em{ color:#043567}.w_weather .more li b{background:url(http://i.tq121.com.cn/i/weather2015/city/iconalls.png) no-repeat scroll -39px -649px;font-size:0;cursor:pointer;float:left;height:12px;width:13px;position:absolute;right:10px;top:14px}.w_weather .more li.add{color:#076ea8;text-align:center;height:15px;width:100%;display:block; border:none;}.w_weather .more li.add a{display:block;font-size:22px;padding:0;width:100%;color:#043567;line-height: 32px;}.city_guonei p{height:25px;line-height:25px;border-bottom:1px dotted #97a5b5;width:410px;margin:0 auto;padding:0 10px}.w_weather .more li:hover{background:#f5fafe}.city_guonei p a{color:#6f6f6f;float:none;font-size:12px;margin-bottom:0;padding:0;text-align:left}.city_guonei p span{color:#d0d0d0;float:right}.w_weather .more li em.w_yj{width:25px;background-position:-177px -476px}.footer{border-top:1px solid #eee7e7;background:#f2f2f8;width:100%;clear:both;overflow:hidden}.footer .block{margin:0 auto 10px;position:relative;width:1000px;padding:20px 0 0;margin-bottom:0;height:auto;background:#f2f2f8}.footer h2{width:145px;height:38px;background:url(../i/index_icons.png) no-repeat -160px 3px;float:left;margin:10px 3px 3px 0;overflow:hidden}.footer h2 a{display:block;width:145px;height:38px;text-indent:-9999px}.footer .Lcontent{width:592px;height:130px;float:left}.footer dl{width:115px;padding:15px;float:left;margin:0 3px 3px 0;display:inline}.footer dl dt{line-height:30px;padding-bottom:5px}.footer dl dt h3,.friendLink h3{font-weight:400;font-size:20px}.footer dl dd p,.friendLink p{line-height:25px}.footer dl p a,.friendLink p a{padding-right:10px;color:#959595;font-size:14px}.footer dl p a:hover,.friendLink p a:hover{color:#ee842f}.friendLink{width:390px;float:right;padding:15px 0 13px 0}.friendLink h3{padding-bottom:5px;line-height:30px}.footer .last{border-top:1px solid #ebebeb;margin-top:5px;padding-top:5px}.serviceinfo{border-top:1px dashed #d5d4d4;padding:20px 0 20px 15px;font-size:12px;line-height:20px;clear:both}.serviceinfo p{height:20px}.serviceinfo p span{display:block;width:334px;float:left}.serviceinfo a,.serviceinfo b{color:#076ea8;font-weight:400}.serviceinfo b{color:#252525}.aboutUs{background-color:#7d7d7d;color:#fff;height:45px;line-height:45px;text-align:center;font-size:12px;width:100%;margin:0 auto}.aboutUs a{color:#fff}.search-box #show li{font-size:15px;height:28px;line-height:28px;list-style-type:none;margin:0;overflow:hidden;padding:0}.search-box #show ul li{cursor:pointer;text-indent:10px}.city_guonei dl dd.jind a{width:61px;overflow:hidden;height:20px}#show ul{border:1px solid #c2d0e7}#show ul li b{color:#f60;font-weight:700}#show ul .select{background-color:#f2f2f8;text-align:left;margin:0;padding:0;color:#252525}#show ul .unselect{padding:0;margin:0}.search-box #show ul{border:1px solid #c2d0e7}#show{display:none}#show{background-color:#FFF;color:#000;display:none;margin:0;overflow-x:hidden;overflow-y:auto;padding:0;position:absolute;top:45px;left:0;width:374px;z-index:99}#show ul{border:1px solid #C2D0E7}#show ul span{background:url(http://i.tq121.com.cn/i/weather2015/index/smile.jpg) no-repeat 10px 5px;display:block;font-size:14px;height:34px;line-height:34px;text-indent:38px}#show li{overflow:hidden;font-size:15px;height:28px;line-height:28px;list-style-type:none;margin:0;padding:0;text-indent:10px}#show ul li b{color:#ee842f;font-weight:normal}.select{background-color:#f2f2f8;margin:0;padding:0;text-align:left}.unselect{margin:0;padding:0}big.icon{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue22.png);height:30px;width:23px;margin:0 2px;position:relative;top:8px;}big.ic{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue22-1.png);height:30px;margin-top:9px;width:28px;float:left}.around span.move,.greatEvent p.time,.livezs li i,.right .list li,.rollLeft,.rollRight,ul.botIcon li{background-image:url(http://i.tq121.com.cn/i/weather2015/city/iconalls.png);background-repeat:no-repeat;display:block}big{margin:0 auto;background-repeat:no-repeat;background-position:-640px 240px;display:block}big.jpg30{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue30.png);height:30px;width:31px}big.png30{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue30.png);height:30px;width:31px}big.jpg50{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue50.png);height:50px;width:50px}big.jpg80{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue80.png);height:80px;width:80px}big.png40{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue30.png);height:30px;width:30px}big.png80{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue80.png);height:80px;width:80px}.around li:hover big.jpg30{background-image:url(http://i.tq121.com.cn/i/weather2015/png/white30.png);height:30px;width:31px}.on big.jpg50{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue50.png);height:50px;width:50px}.on big.jpg80{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue80.png);height:80px;width:80px}.sk big.jpg80{background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue80.png);height:80px;width:80px}big.d0{background-position:0 0}big.d1{background-position:-80px 0}big.d2{background-position:-160px 0}big.d3{background-position:-240px 0}big.d4{background-position:-320px 0}big.d5{background-position:-400px 0}big.d6{background-position:-480px 0}big.d7{background-position:-560px 0}big.d8{background-position:-640px 0}big.d9{background-position:0 -80px}big.d00{background-position:0 0}big.d01{background-position:-80px 0}big.d02{background-position:-160px 0}big.d03{background-position:-240px 0}big.d04{background-position:-320px 0}big.d05{background-position:-400px 0}big.d06{background-position:-480px 0}big.d07{background-position:-560px 0}big.d08{background-position:-640px 0}big.d09{background-position:0 -80px}big.d10{background-position:-80px -80px}big.d11{background-position:-160px -80px}big.d12{background-position:-240px -80px}big.d13{background-position:-320px -80px}big.d14{background-position:-400px -80px}big.d15{background-position:-480px -80px}big.d16{background-position:-560px -80px}big.d17{background-position:-640px -80px}big.d18{background-position:0 -160px}big.d19{background-position:-80px -160px}big.d20{background-position:-160px -160px}big.d21{background-position:-240px -160px}big.d22{background-position:-320px -160px}big.d23{background-position:-400px -160px}big.d24{background-position:-480px -160px}big.d25{background-position:-560px -160px}big.d26{background-position:-640px -160px}big.d27{background-position:0 -240px}big.d28{background-position:-80px -240px}big.d29{background-position:-160px -240px}big.d30{background-position:-240px -240px}big.d31{background-position:-320px -240px}big.d32{background-position:-400px -240px}big.d33{background-position:-480px -240px}big.d53{background-position:-560px -240px}big.d57{background-position:-720px 0}big.d32{background-position:-720px -80px}big.d49{background-position:-720px -160px}big.d58{background-position:-720px -240px}big.d54{background-position:-800px 0}big.d55{background-position:-800px -80px}big.d56{background-position:-800px -160px}big.d301{background-position:-880px 0}big.d302{background-position:-880px -80px}big.n0{background-position:0 -320px}big.n1{background-position:-80px -320px}big.n2{background-position:-160px -320px}big.n3{background-position:-240px -320px}big.n4{background-position:-320px -320px}big.n5{background-position:-400px -320px}big.n6{background-position:-480px -320px}big.n7{background-position:-560px -320px}big.n8{background-position:-640px -320px}big.n9{background-position:0 -400px}big.n00{background-position:0 -320px}big.n01{background-position:-80px -320px}big.n02{background-position:-160px -320px}big.n03{background-position:-240px -320px}big.n04{background-position:-320px -320px}big.n05{background-position:-400px -320px}big.n06{background-position:-480px -320px}big.n07{background-position:-560px -320px}big.n08{background-position:-640px -320px}big.n09{background-position:0 -400px}big.n10{background-position:-80px -400px}big.n11{background-position:-160px -400px}big.n12{background-position:-240px -400px}big.n13{background-position:-320px -400px}big.n14{background-position:-400px -400px}big.n15{background-position:-480px -400px}big.n16{background-position:-560px -400px}big.n17{background-position:-640px -400px}big.n18{background-position:0 -480px}big.n19{background-position:-80px -480px}big.n20{background-position:-160px -480px}big.n21{background-position:-240px -480px}big.n22{background-position:-320px -480px}big.n23{background-position:-400px -480px}big.n24{background-position:-480px -480px}big.n25{background-position:-560px -480px}big.n26{background-position:-640px -480px}big.n27{background-position:0 -560px}big.n28{background-position:-80px -560px}big.n29{background-position:-160px -560px}big.n30{background-position:-240px -560px}big.n31{background-position:-320px -560px}big.n32{background-position:-400px -560px}big.n33{background-position:-480px -560px}big.n53{background-position:-560px -560px}big.n57{background-position:-720px -320px}big.n32{background-position:-720px -400px}big.n49{background-position:-720px -480px}big.n58{background-position:-720px -560px}big.n54{background-position:-800px -320px}big.n55{background-position:-800px -400px}big.n56{background-position:-800px -480px}big.n301{background-position:-880px -320px}big.n302{background-position:-880px -400px}.adposter_pos{background:transparent url(#) repeat scroll 0 0;left:-10000px;position:absolute;z-index:9}.ad{height:93px;overflow:hidden}.ad .ad1{height:90px;overflow:hidden;width:260px}.ad .ad2{float:right;height:100%;margin-left:10px;width:728px}.ad4{height:250px;margin-top:15px;width:300px}.ad3{height:94px;margin-top:15px;width:680px}#adposter_6287{z-index:1006}div#ab_yjfk{left:50%;margin-left:500px;position:fixed;top:370px}.topad_bg{}#duilian{z-index:1006}.nav a:hover{color:#ee842f}#ab_yjjy{margin-right:500px;position:fixed;right:50%;top:395px;}#abs { z-index:10 ;display:none;}.provinceLinks { position:absolute; z-index:3000; top:27px; left:0px; width:100%; height:175px; background-color:#f6fcff;  display:none;box-shadow: -7px 13px 16px -23px #000;}.w_weather em.s:hover{color:#ee842f;}
.provinceLinks dl { float:left; padding:10px 0px; }
.provinceLinks dl dt { width:46px; height:40px; background:url(http://i.tq121.com.cn/i/weather2015/zt/t.png) no-repeat -1px -1px; font-size:16px; line-height:20px; padding:13px 10px; text-align:center; float:left; margin-right:4px; }
.provinceLinks dl .last { padding:13px 5px; width:56px; }
.provinceLinks a { color:#076ea8; }
.provinceLinks a:hover { color:#ee842f; }

.provinceLinks dl dd { font-size:16px; float:left; padding:5px 0px; height:56px; line-height:28px; }
.provinceLinks dl dd a { padding:0px 19px;font-size:16px;font-family:"Microsoft Yahei",Tahoma,SimSun; }
.provinceLinks .midBlock dl dt a { color:#fff; font-weight: bold;font-size:16px; }
.sjz{ background:url(http://i.tq121.com.cn/i/weather2015/zt/t.png) no-repeat scroll -9px -72px}
.provinceLinks .line { border-bottom:0; height:0px; overflow:hidden; float:left; width:1000px; }
.nav_li { background: #043567; height: 30px; width: 100%; min-width: 1000px; line-height: 30px; }
.nav_li_box { position: relative; width: 1000px; z-index: 1000; margin: 0 auto; }
.nav_li_left { float: left; color: #fff; font-size: 12px; width: 720px; }
.nav_li_left a { color: #fff; display: block; float: left; margin-right: 10px; font-size: 12px; }
.nav_li_right { float: right; width: 280px; text-align: right; }
.nav_li_left span { display: block; float: left; margin-right: 10px; }
.nav_li_right a { color: #fff; font-size: 12px; }
.weather_li_right a.login_li {display:none; background: url(http://i.tq121.com.cn/i/weather2015/index/indexImgs.png) no-repeat -580px -115px; padding-left: 29px; }
.nav_addr { width: 1000px; height: 65px; border-bottom: 1px solid #dee8f2; margin: 0 auto; }
.nav_addr ul li { padding: 0px 8px; margin-bottom: 5px; float: left; }
.nav_addr ul li a { font-size: 14px; }
.nav_addr ul { width: 265px; overflow: hidden; float: left; border-right: 1px solid #e0e9f3; }
.input-btn { width: 20px; right: 10px; top: 10px; position: absolute; height: 20px; background: url("http://i.tq121.com.cn/i/weather2015/index/indexImgs.png") no-repeat -583px -198px; }
.input-btn input#btnZip { color: #fff; cursor: pointer; font-family: 微软雅黑; font-size: 16px; height: 20px; margin-left: 3px; text-indent: 0; width: 20px; background: none; border:none;}
.input-btn input#btnZip:hover { border: none; }
.search { width: 370px; height: 40px; border: 1px solid #cad1d8;border-radius:4px; }
.select_li { display:none;border: 1px solid #a7b5c2; float: left; height: 40px; z-index: 12; left: -1px; overflow: hidden; background: url(http://i.tq121.com.cn/i/weather2015/index/indexImgs.png) no-repeat -544px -143px #fff; position: absolute; top: -1px; width: 60px; }
.select_li p { height: 40px; text-align: center; width: 60px; line-height: 40px; cursor: pointer; }
.select_li b { display: block; background: #fff; height: 40px; line-height: 40px; text-align: center; cursor: pointer; font-weight: normal; }
.select_li b:hover { color: #ee842f }
.select_li b.m_li { background: #eaf4fe }
.search input#txtZip { margin-left: 0px; width: 329px; border: none; background: #fff; height: 40px; border: 0px; }
/*.search input#txtZip { margin-left: 60px; width: 269px; border: none; background: #fff; height: 40px; border-left: 0px solid #ccc; }*/
.search input#txtZip:hover { border: none; border-left: 0px solid #a7b5c2; }
.search-box { margin: 0 0 0 153px; position: relative; width: 373px; }
/*预报详情页*/
.weather_li { background: #043567; height: 40px; width: 100%; min-width: 1000px; line-height: 40px;position:relative;z-index:1501; }
.weather_li_left { float: left; color: #fff; font-size: 12px; width: 620px; }
.weather_li_left a { color: #fff; display: block; float: left; margin-right: 17px; font-size: 12px; }
.weather_li_head { height: 75px; width: 100%; min-width: 1000px; background: #f2f2f8;margin-bottom:10px;position:relative;z-index:1500; }
.weather_li_box { margin: 0 auto; width: 1000px; }
.weather_li_right { float: right; width: 320px; height: 40px; text-align: right; }
.weather_li_right em { font-size: 12px; color: #fff; }
.weather_li_right .w_weather em.s { color: #fff; font-size: 12px; }
.weather_li_right .w_weather a:hover em.s{ color: #ee842f;  }

.weather_li_right .w_weather { margin-top: 0; height: 40px; float: right; }
.login_li { color: #fff; font-size: 12px; }
.w_li_logo { width: 300px; }
.w_li_logo a { background:url(http://i.tq121.com.cn/i/weather2015/index/indexImgs.png) no-repeat -335px -213px; display: block; height: 34px; margin-top: 20px; float: left; border-right: 1px solid #ccc; width: 162px; }
.w_li_logo span { float: left; margin-top: 20px; display: inline;  margin-left: 15px; font-size: 22px; color:#252525}
.weather_li_left div.more_li { padding: 0px 25px; padding-left:15px;position: relative; float: left; cursor: pointer; background: url("http://i.tq121.com.cn/i/weather2015/index/indexImgs.png") no-repeat -548px -144px; display: block; font-size:12px;}
.weather_li_left a:hover { color: #ee842f; }
.weather_li_left div.more_li:hover { color: #ee842f; background:url("http://i.tq121.com.cn/i/weather2015/index/indexImgs.png") no-repeat #fff -548px -277px }
.weather_li_open { width: 425px; height: 100px; position: absolute; display: none; left: -323px; top: 40px; border: 1px solid #d8d8da; border-top: none; background: #fff; }
.weather_li_open p { border-bottom: 1px dashed #d7d8dd;height: 50px;line-height: 50px; margin: 0 auto 0 3%;  width: 93%;
}
.weather_li_open p a { color: #252525; margin-right: 0; padding: 0px 9px; font-size: 12px; }
.weather_li_open p.erp a{padding:0px 8px;}
.weather_li_open p a:hover { color: #ee842f; }
.w_weather a.add {line-height:100px;overflow:hidden; }
big.icon { background-image: url(http://i.tq121.com.cn/i/weather2015/png/white22.png); }
.w_weather span em { font-size: 12px; }
.w_weather span { height: 40px; line-height: 40px; }
.w_weather em.s { height: 40px; line-height: 40px; }
#w_weather a.add{
    background-position: -12px -616px;
}
#w_weather a.dz_right {
    background-position: -60px -615px;
}
#w_weather a.dz_down{
    background-position: -38px -615px;
}
.locationSearch .search{border:none;}
.w_weather a:hover big.icon{
    background-image:url(http://i.tq121.com.cn/i/weather2015/png/blue22-2.png);
    }
.w_weather a:hover em{
color:#ee842f}
.footer .serviceinfo a{font-size:12px;}
input::-webkit-search-cancel-button{
display: none;
}
input::-ms-clear{
display: none;
}
</style>

<div class="weather_li">
	<div class="nav_li_box">
    	<div class="weather_li_left">
        	<a target="_blank" href="http://www.weather.com.cn/">首页</a>
            <a target="_blank" href="http://www.weather.com.cn/forecast/">预报</a>
            <a target="_blank" href="http://www.weather.com.cn/satellite/">云图</a>
            <a target="_blank" href="http://www.weather.com.cn/radar/">雷达</a>
            <a target="_blank" href="http://www.weather.com.cn/live/">临近预报</a>
            <a target="_blank" href="http://products.weather.com.cn/">专业产品</a>
            <a target="_blank" href="http://news.weather.com.cn/">资讯</a>
            <a target="_blank" href="http://www.weather.com.cn/life/">生活</a>
            <a href="http://www.weather.com.cn/forecast/skiweather.shtml" target="_blank">滑雪</a>
            <a target="_blank" href="http://www.weatherdt.com/">产创平台</a>
            <a target="_blank" href="" class="shengjz"></a>
            <div href="javascript:void(0)" class="more_li">更多
            <div class="weather_li_open">
        		<p>
                    <a target="_blank" href="http://www.weather.com.cn/alarm/">预警</a>
                    <a target="_blank" href="http://typhoon.weather.com.cn/">台风路径</a>
                    <a target="_blank" href="http://www.weather.com.cn/space/">空间天气</a>
                    <a target="_blank" href="http://p.weather.com.cn/">图片</a>
                    <a target="_blank" href="http://www.weather.com.cn/video/">视频</a>
                    <a target="_blank" href="http://www.weather.com.cn/zt/">专题</a>
                    <a target="_blank" href="http://www.weather.com.cn/air/">环境</a>
                    <a  href="http://www.weather.com.cn/science/" target="_blank">科普</a>
                </p>
                <p style="border:none;" class="erp">
                    <a target="_blank" href="http://www.weather.com.cn/trip/">旅游</a>
                    <a target="_blank" href="http://www.sportsweather.cn/golf/">高尔夫</a>
                    <a target="_blank" href="http://www.weather.com.cn/jt/">交通</a>
                    <a target="_blank" href="http://www.weather.com.cn/fzjz/">减灾</a>
                    <a target="_blank" href="http://www.weather.com.cn/climate/">气候变化</a>
                    <a target="_blank" href="http://marketing.weather.com.cn/">商业合作</a>
                    <a target="_blank" href="http://www.weather.com.cn/province/">省级站</a>
                    <a target="_blank" href="http://club.weather.com.cn/">社区</a>
                </p>
            </div>

            </div>
        </div>
        <div class="weather_li_right">
        	  <div class="w_weather" id="w_weather">	</div><a href="#" class="login_li">登录</a>
        </div>

    </div>
</div>
<div class="weather_li_head">
	<div class="weather_li_box">
    	<div class="w_li_logo fl">
        	<a href="http://www.weather.com.cn/"></a>
            <span></span>
        </div>
        <div class="search-box fl" style=" float:right; margin-top:16px;">
<div class="search clearfix">
<div class="select_li">
	<p>天气</p>
    <b class="m_li">天气</b>
    <b>资讯</b>
</div>
<input type="text" value="输入城市名、景点名 查天气" id="txtZip" class="textinput text fl">

<span class="input-btn"><input type="button" value="" id="btnZip" class="btn ss fl"></span>
<div class="clear"></div>
</div>
<div class="inforesult"> </div>
<div id="show">
<ul>
</ul>
</div>
<div class="city-box">
<div class="city-tt"> <a href="javascript:void(0)" class="cur">正在热搜</a> <a href="javascript:void(0)" >本地周边</a> <b></b> </div>
<div class="w_city city_guonei" style="display:block; padding-bottom:0">
<dl>
<dt>国内</dt>
<dd>  <a href="http://www.weather.com.cn/weather1d/101010100.shtml#search">北京</a> <a href="http://www.weather.com.cn/weather1d/101020100.shtml#search">上海</a> <a href="http://www.weather.com.cn/weather1d/101210101.shtml#search">杭州</a> <a href="http://www.weather.com.cn/weather1d/101280101.shtml#search">广州</a> <a href="http://www.weather.com.cn/weather1d/101200101.shtml#search">武汉</a> <a href="http://www.weather.com.cn/weather1d/101190101.shtml#search">南京</a> <a href="http://www.weather.com.cn/weather1d/101280601.shtml#search">深圳</a> <a href="http://www.weather.com.cn/weather1d/101190401.shtml#search">苏州</a> <a href="http://www.weather.com.cn/weather1d/101230201.shtml#search">厦门</a> <a href="http://www.weather.com.cn/weather1d/101220101.shtml#search">合肥</a> <a href="http://www.weather.com.cn/weather1d/101250101.shtml#search">长沙</a>  <a href="http://www.weather.com.cn/weather1d/101270101.shtml#search">成都</a>
</dd>
</dl>
<dl>
<dt>国际</dt>
<dd> <a href="http://www.weather.com.cn/weather1d/102010100.shtml#search">首尔</a> <a href="http://www.weather.com.cn/weather1d/104010100.shtml#search">新加坡</a> <a href="http://www.weather.com.cn/weather1d/106010100.shtml#search">曼谷</a> <a href="http://www.weather.com.cn/weather1d/401110101.shtml#search">纽约</a> <a href="http://www.weather.com.cn/weather1d/124020100.shtml#search">迪拜</a> <a href="http://www.weather.com.cn/weather1d/103163100.shtml#search">大阪</a> <a href="http://www.weather.com.cn/weather1d/601020101.shtml#search">悉尼</a> <a href="http://www.weather.com.cn/weather1d/601060101.shtml#search">墨尔本</a> <a href="http://www.weather.com.cn/weather1d/401040101.shtml#search">洛杉矶</a> <a href="http://www.weather.com.cn/weather1d/105010100.shtml#search">吉隆坡</a> </dd>
</dl>
<dl>
<dt>景点</dt>
<dd> <a href="http://www.weather.com.cn/weather1d/10101010018A.shtml#search">故宫</a> <a href="http://www.weather.com.cn/weather1d/10130051008A.shtml#search">阳朔漓江</a> <a href="http://www.weather.com.cn/weather1d/10118090107A.shtml#search">龙门石窟</a> <a href="http://www.weather.com.cn/weather1d/10109022201A.shtml#search">野三坡</a> <a href="http://www.weather.com.cn/weather1d/10101020015A.shtml#search">颐和园</a> <a href="http://www.weather.com.cn/weather1d/10127190601A.shtml#search">九寨沟</a> <a href="http://www.weather.com.cn/weather1d/10102010007A.shtml#search">东方明珠</a> <a href="http://www.weather.com.cn/weather1d/10125150503A.shtml#search">凤凰古城</a> <a href="http://www.weather.com.cn/weather1d/10111010119A.shtml#search">秦始皇陵</a> <a href="http://www.weather.com.cn/weather1d/10125060301A.shtml#search">桃花源</a> </dd>
</dl>
<dl style="margin-bottom:5px;border:none;">
<dt>高球</dt>
<dd>
<a href="http://www.sportsweather.cn/weather/10102090003F.shtml#search">佘山</a>
<a href="http://www.sportsweather.cn/weather/10129010601F.shtml#search">春城湖畔</a>
 <a href="http://www.sportsweather.cn/weather/10101070004F.shtml#search">华彬庄园</a>
 <a href="http://www.sportsweather.cn/weather/10128060113F.shtml#search">观澜湖</a>
<a href="http://www.sportsweather.cn/weather/10131010107F.shtml#search">依必朗</a>
 <a href="http://www.sportsweather.cn/weather/10102080001F.shtml#search">旭宝</a>
 <a href="http://www.sportsweather.cn/weather/10131021101F.shtml#search">博鳌</a>
<a href="http://www.sportsweather.cn/weather/10129140501F.shtml#search">玉龙雪山</a>
<a href="http://www.sportsweather.cn/weather/10128010103F.shtml#search">番禺南沙</a>
 <a href="http://www.sportsweather.cn/weather/10101040001F.shtml#search">东方明珠</a>
</dd></dl>
</div>
<div class="w_city city_guonei gn">
<dl>
<dt>地区</dt>
<dd class="diq"></dd>
</dl>
<dl style="border:none;margin-bottom:5px;">
<dt>景点</dt>
<dd class="jind"></dd>
</dl>
</div>
</div>
</div>
    </div>
</div>
<script type="text/javascript" src="http://i.tq121.com.cn/j/core.js"></script>

<div class="topad_bg">
<div class="box" style="padding:0;">
				<!--顶通两个begin-->
				<div class="ad clearfix post_st">
					<div id="adposter_6125" class="ad1 fl" style="float:right;width:260px; height:90px; overflow:hidden;">
						<script>
						(function() {
						    var s = "_" + Math.random().toString(36).slice(2);
						    document.write('<div id="' + s + '"></div>');
						    (window.slotbydup=window.slotbydup || []).push({
						        id: '3011945',
						        container: s,
						        size: '260,90',
						        display: 'inlay-fix'
						    });
						})();
						</script>
					</div>
					<div id="adposter_6126" class="ad2 fl" style="float:left;margin-left:0;">
							<script>
							(function() {
							    var s = "_" + Math.random().toString(36).slice(2);
							    document.write('<div id="' + s + '"></div>');
							    (window.slotbydup=window.slotbydup || []).push({
							        id: '3011939',
							        container: s,
							        size: '728,90',
							        display: 'inlay-fix'
							    });
							})();
							</script>
					</div>
				</div>
			<!--顶通两个end-->
		</div>
		</div>

<link href="http://i.tq121.com.cn/c/weather2015/common.css" rel="stylesheet" type="text/css">
<link href="http://i.tq121.com.cn/c/weather2015/bluesky/c_1d.css" rel="stylesheet" type="text/css">


<input id="whichDay" type="hidden" value="today" />

<div class="con today clearfix">


	<div class="left fl">
		<div class="ctop clearfix">
			<div class="crumbs fl">
				<a href="http://bj.weather.com.cn" target="_blank">北京</a>
				<span>></span>
				 <span>城区</span>
			</div>
			<div class="time fr"></div>
		</div>
		<ul id="someDayNav" class="clearfix cnav">
			<li class="on">
				<a href="/weather1d/101010100.shtml">今天</a>
			</li>
			<li class="hover">
				<a href="/weather/101010100.shtml">7天</a>
			</li>
			<li>
				<a href="/weather15d/101010100.shtml">更长时间</a>
			</li>
						<li>
				<a href="/weather40d/101010100.shtml">天气日历</a><span></span>
			</li>
			<li>
				<a href="http://products.weather.com.cn/product/radar/index/procode/JC_RADAR_AZ9010_JB" target="_blank">雷达图</a>
			</li>
		</ul>
		<div class="today clearfix" id="today">
					<input type="hidden" id="hidden_title" value="11月27日12时 周日  晴  8/-4°C" />
<input type="hidden" id="update_time" value="11:30"/>
<input type="hidden" id="fc_24h_internal_update_time" value="2016-11-27 11:30"/>
<div class="t">
<div class="sk">
<div class="zs limit">
<i></i><span>限行</span>
<em>不限行</em>
</div>
</div>
<ul class="clearfix">
<li>
<h1>27日白天</h1>
<big class="jpg80 d00"></big>
<p class="wea" title="晴">晴</p>
<div class="sky">
<span class="txt lv1">天空蔚蓝</span>
<i class="icon"></i>
<div class="skypop">
<h3>蓝天预报综合天气现象、能见度、空气质量等因子，预测未来一周的天空状况。</h3>
<ul>
<li class="lv1">
<em></em><span>天空蔚蓝</span>
<b>可见透彻蓝天，或有蓝天白云美景。</b>
</li>
<li class="lv2">
<em></em><span>天空淡蓝</span>
<b>天空不够清澈，以浅蓝色为主。</b>
</li>
<li class="lv3">
<em></em><span>天空阴沉</span>
<b>阴天或有雨雪，天空灰暗。</b>
</li>
<li class="lv4">
<em></em><span>天空灰霾</span>
<b>出现霾或沙尘，天空灰蒙浑浊。</b>
</li>
</ul>
<i class="s"></i>
</div>
</div>
<p class="tem">
<span>8</span><em>°C</em>
</p>
<p class="win">
<i class="N"></i>
<span class="" title="北风">3-4级</span>
</p>
<p class="sun sunUp"><i></i>
<span>日出 07:12</span>
</p>
<div class="slid"></div>
</li>
<li>
<h1>27日夜间</h1>
<big class="jpg80 n00"></big>
<p class="wea" title="晴">晴</p>
<div class="sky">
</div>
<p class="tem">
<span>-4</span><em>°C</em>
</p>
<p class="win"><i class=""></i><span class="" title="无持续风向">微风</span></p>
<p class="sun sunDown"><i></i>
<span>日落 16:51</span>
</p>
</li>
</ul>
</div>
							<input type="hidden" id="fc_3h_internal_update_time" value="2016-11-27 11:30"/>
<div class="curve_livezs" id="curve">
<div class="time">
</div>
<div class="wpic">
</div>
<div id="biggt" class="biggt">
</div>
<div class="tem">
</div>
<div class="winf">
</div>
<div class="winl">
</div>
</div>
<script>
var hour3data={"1d":["27日08时,d01,多云,2℃,北风,3-4级,2","27日11时,d00,晴,6℃,北风,微风,1","27日14时,d00,晴,7℃,北风,微风,1","27日17时,d00,晴,2℃,北风,微风,1","27日20时,n00,晴,0℃,北风,微风,0","27日23时,n00,晴,-1℃,无持续风向,微风,0","28日02时,n00,晴,-3℃,无持续风向,微风,0","28日05时,n00,晴,-3℃,无持续风向,微风,0","28日08时,d00,晴,0℃,无持续风向,微风,2"],"23d":[["27日08时,d01,多云,2℃,北风,3-4级,2","27日11时,d00,晴,6℃,北风,微风,1","27日14时,d00,晴,7℃,北风,微风,1","27日17时,d00,晴,2℃,北风,微风,1","27日20时,n00,晴,0℃,北风,微风,0","27日23时,n00,晴,-1℃,无持续风向,微风,0","28日02时,n00,晴,-3℃,无持续风向,微风,0","28日05时,n00,晴,-3℃,无持续风向,微风,0"],["03日08时,d00,晴,0℃,无持续风向,微风,1","03日14时,d00,晴,7℃,无持续风向,微风,1","03日20时,n00,晴,0℃,无持续风向,微风,0","04日02时,n02,阴,-1℃,无持续风向,微风,0"]],"7d":[["27日08时,d01,多云,2℃,北风,3-4级,2","27日11时,d00,晴,6℃,北风,微风,1","27日14时,d00,晴,7℃,北风,微风,1","27日17时,d00,晴,2℃,北风,微风,1","27日20时,n00,晴,0℃,北风,微风,0","27日23时,n00,晴,-1℃,无持续风向,微风,0","28日02时,n00,晴,-3℃,无持续风向,微风,0","28日05时,n00,晴,-3℃,无持续风向,微风,0"],["28日08时,d00,晴,0℃,无持续风向,微风,2","28日11时,d00,晴,4℃,无持续风向,微风,2","28日14时,d00,晴,6℃,无持续风向,微风,2","28日17时,d00,晴,1℃,无持续风向,微风,2","28日20时,n00,晴,0℃,无持续风向,微风,0","28日23时,n53,霾,0℃,无持续风向,微风,0","29日02时,n53,霾,-1℃,无持续风向,微风,0","29日05时,n53,霾,-2℃,无持续风向,微风,0"],["29日08时,d53,霾,-3℃,无持续风向,微风,4","29日11时,d53,霾,1℃,无持续风向,微风,4","29日14时,d53,霾,2℃,无持续风向,微风,4","29日17时,d53,霾,0℃,无持续风向,微风,4","29日20时,n53,霾,0℃,无持续风向,微风,0","29日23时,n53,霾,-1℃,无持续风向,微风,0","30日02时,n53,霾,-1℃,无持续风向,微风,0","30日05时,n53,霾,-3℃,无持续风向,微风,0"],["30日08时,d53,霾,-1℃,无持续风向,微风,4","30日14时,d53,霾,6℃,无持续风向,微风,4","30日20时,n53,霾,2℃,无持续风向,微风,0","01日02时,n00,晴,0℃,无持续风向,微风,0"],["01日08时,d00,晴,0℃,无持续风向,微风,2","01日14时,d00,晴,6℃,无持续风向,微风,2","01日20时,n00,晴,-1℃,无持续风向,微风,0","02日02时,n00,晴,-2℃,无持续风向,微风,0"],["02日08时,d00,晴,-1℃,无持续风向,微风,1","02日14时,d00,晴,6℃,无持续风向,微风,1","02日20时,n00,晴,0℃,无持续风向,微风,0","03日02时,n00,晴,0℃,无持续风向,微风,0"],["03日08时,d00,晴,0℃,无持续风向,微风,1","03日14时,d00,晴,7℃,无持续风向,微风,1","03日20时,n00,晴,0℃,无持续风向,微风,0","04日02时,n02,阴,-1℃,无持续风向,微风,0"]]}
</script>

		</div>
		<!--条形广告位begin-->
			<div class="ad3" id="adposter_6122">
				<script>
				(function() {
				    var s = "_" + Math.random().toString(36).slice(2);
				    document.write('<div id="' + s + '"></div>');
				    (window.slotbydup=window.slotbydup || []).push({
				        id: '3011953',
				        container: s,
				        size: '680,90',
				        display: 'inlay-fix'
				    });
				})();
				</script>
			</div>
			<!--条形广告位end-->
				<script>
var observe24h_data = {"od":{"od0":"20161127130000","od1":"北京_南郊观象台","od2":[{"od21":"13","od22":"7","od23":"354","od24":"北风","od25":"3","od26":"0","od27":"18","od28":""},{"od21":"12","od22":"6","od23":"345","od24":"北风","od25":"3","od26":"0","od27":"18","od28":"222"},{"od21":"11","od22":"6","od23":"11","od24":"北风","od25":"3","od26":"0","od27":"19","od28":"232"},{"od21":"10","od22":"6","od23":"19","od24":"北风","od25":"2","od26":"0","od27":"21","od28":"243"},{"od21":"09","od22":"5","od23":"31","od24":"东北风","od25":"2","od26":"0","od27":"23","od28":"255"},{"od21":"08","od22":"4","od23":"24","od24":"东北风","od25":"2","od26":"0","od27":"26","od28":"265"},{"od21":"07","od22":"4","od23":"17","od24":"北风","od25":"2","od26":"0","od27":"26","od28":"275"},{"od21":"06","od22":"5","od23":"9","od24":"北风","od25":"1","od26":"0","od27":"26","od28":"284"},{"od21":"05","od22":"6","od23":"323","od24":"西北风","od25":"3","od26":"0","od27":"23","od28":"293"},{"od21":"04","od22":"6","od23":"313","od24":"西北风","od25":"3","od26":"0","od27":"26","od28":"301"},{"od21":"03","od22":"4","od23":"325","od24":"西北风","od25":"2","od26":"0","od27":"39","od28":"307"},{"od21":"02","od22":"-1","od23":"333","od24":"西北风","od25":"1","od26":"0","od27":"88","od28":"313"},{"od21":"01","od22":"0","od23":"211","od24":"西南风","od25":"1","od26":"0","od27":"83","od28":"314"},{"od21":"00","od22":"1","od23":"182","od24":"南风","od25":"1","od26":"0","od27":"80","od28":"315"},{"od21":"23","od22":"1","od23":"111","od24":"东风","od25":"1","od26":"0","od27":"76","od28":"316"},{"od21":"22","od22":"2","od23":"231","od24":"西南风","od25":"1","od26":"0","od27":"75","od28":"316"},{"od21":"21","od22":"2","od23":"0","od24":"北风","od25":"0","od26":"0","od27":"78","od28":"316"},{"od21":"20","od22":"2","od23":"236","od24":"西南风","od25":"1","od26":"0","od27":"73","od28":"315"},{"od21":"19","od22":"3","od23":"203","od24":"西南风","od25":"1","od26":"0","od27":"71","od28":"313"},{"od21":"18","od22":"3","od23":"266","od24":"西风","od25":"1","od26":"0","od27":"66","od28":"311"},{"od21":"17","od22":"5","od23":"197","od24":"南风","od25":"1","od26":"0","od27":"59","od28":"308"},{"od21":"16","od22":"6","od23":"214","od24":"西南风","od25":"2","od26":"0","od27":"54","od28":"304"},{"od21":"15","od22":"7","od23":"225","od24":"西南风","od25":"2","od26":"0","od27":"47","od28":"302"},{"od21":"14","od22":"7","od23":"209","od24":"西南风","od25":"2","od26":"0","od27":"41","od28":"298"},{"od21":"13","od22":"6","od23":"215","od24":"西南风","od25":"2","od26":"0","od27":"50","od28":"294"}]}};
</script>
<div id="weatherChart">
</div>
		<div class="clearfix" id="weatherChart"></div>
				<!--条形广告位begin-->
			<div class="ad3" id="adposter_6298">
				<div class="ad3" style="background:url(http://i.tq121.com.cn/i/loading.gif) center center no-repeat #eee"></div>
			</div>
			<!--条形广告位end-->
		<div class="livezs">
			<div class="t clearfix">
				<h1>生活指数</h1>
			</div>
						<script src="http://i.tq121.com.cn/j/weather2015/pagefilp.js" type=text/javascript></script>
<style>
.pageflip {right: 0px; float: right; position: relative; top: 0px}
 .pageflip IMG {z-index: 99; right: -3px; width: 30px; position: absolute; top: -1px; height: 30px; ms-interpolation-mode: bicubic}
 .pageflip .msg_block {right: 0px; background: url(http://i.tq121.com.cn/i/weather2015/png/subscribe.png) no-repeat right top; overflow: hidden; width: 25px; position: absolute; top: 0px; height: 25px}
</style>
<input type="hidden" id="zs_7d_update_time" value="2016-11-27 12:00:00.0"/>
<ul class="clearfix">
<li class="li1 hot">
<i></i>
<span>中等</span>
<em>紫外线指数</em>
<p>涂擦SPF大于15、PA+防晒护肤品。</p>
</li>
<li class="li2 hot">
<i></i>
<span>易发</span>
<em>感冒指数</em>
<p>昼夜温差大，易感冒。</p>
</li>
<li class="li3 hot" id="chuanyi">
<a href="http://www.weather.com.cn/forecast/ct.shtml?areaid=101010100">
<div class="pageflip">
<IMG src="http://i.tq121.com.cn/i/weather2015/png/page_flip.png">
<div class=msg_block>
</div>
</div>
<i></i>
<span>较冷</span>
<em>穿衣指数</em>
<p>建议着厚外套加毛衣等服装。</p>
</a>
</li>
<li class="li4 hot">
<i></i>
<span>较适宜</span>
<em>洗车指数</em>
<p>无雨且风力较小，易保持清洁度。</p>
</li>
<li class="li5 hot">
<i></i>
<span>较不宜</span>
<em>运动指数</em>
<p>天气寒冷，推荐您进行室内运动。</p>
</li>
<li class="li6 hot">
<i></i>
<span>良</span>
<em>空气污染扩散指数</em>
<p>气象条件有利于空气污染物扩散。</p>
</li>
</ul>
					</div>
		<!--条形广告位begin-->
			<div class="ad3" id="adposter_6299">
				<div class="ad3" style="background:url(http://i.tq121.com.cn/i/loading.gif) center center no-repeat #eee"></div>
			</div>
			<!--条形广告位end-->
		<div id="around" class="around">
								<div class="aro_city" style="display:block;">
					<input type="hidden" id="around_city_china_update_time" value="2016112708"/>
<h1 class="clearfix city">
<span class="move">周边地区</span>
<em>|</em>
<span>周边景点</span>
<i>2016-11-27 11:30更新</i>
</h1>
<ul class="clearfix city">
<li>
<a href="http://www.weather.com.cn/weather1d/101100101.shtml#around2" target="_blank">
<span>太原</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-6°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090101.shtml#around2" target="_blank">
<span>石家庄</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>11/-2°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101030100.shtml#around2" target="_blank">
<span>天津</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>7/-2°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090301.shtml#around2" target="_blank">
<span>张家口</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>2/-7°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090701.shtml#around2" target="_blank">
<span>沧州</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-4°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090501.shtml#around2" target="_blank">
<span>唐山</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>7/-6°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090201.shtml#around2" target="_blank">
<span>保定</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>10/-4°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090601.shtml#around2" target="_blank">
<span>廊坊</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>9/-5°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090218.shtml#around2" target="_blank">
<span>涿州</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>10/-5°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090312.shtml#around2" target="_blank">
<span>涿鹿</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>3/-5°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090604.shtml#around2" target="_blank">
<span>香河</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>9/-6°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/101090609.shtml#around2" target="_blank">
<span>三河</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>9/-6°C</i>
</a>
</li>
</ul>
				</div>
				<div class="aro_view">
					<input type="hidden" id="around_city_travel_update_time" value="20161127080000"/>
<h1 class="clearfix view">
<span>周边地区</span>
<em>|</em>
<span class="move">周边景点</span>
<i>2016-11-27 11:30更新</i>
</h1>
<ul class="clearfix view">
<li>
<a href="http://www.weather.com.cn/weather1d/10101010002A.shtml#around2" target="_blank">
<span>北京市规划展览馆</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010003A.shtml#around2" target="_blank">
<span>明城墙遗址公园</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010012A.shtml#around2" target="_blank">
<span>景山公园</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010013A.shtml#around2" target="_blank">
<span>什刹海</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010004A.shtml#around2" target="_blank">
<span>天坛公园</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010005A.shtml#around2" target="_blank">
<span>南锣鼓巷</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010006A.shtml#around2" target="_blank">
<span>北京国子监</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010007A.shtml#around2" target="_blank">
<span>北京孔庙</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010008A.shtml#around2" target="_blank">
<span>月坛公园</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010009A.shtml#around2" target="_blank">
<span>中国地质博物馆</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010010A.shtml#around2" target="_blank">
<span>北海公园</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
<li>
<a href="http://www.weather.com.cn/weather1d/10101010011A.shtml#around2" target="_blank">
<span>北京海洋馆</span>
<p class="img clearfix">
<big class="jpg30 d00"></big>
<em>/</em>
<big class="jpg30 n00"></big>
</p>
<i>8/-3°C</i>
</a>
</li>
</ul>
				</div>
						</div>
					 <div class="hdImgs">
				<h3><span><a href="http://p.weather.com.cn/" target="_blank">>></a></span><a href="http://p.weather.com.cn/" target="_blank">高清图集</a></h3>
				<a id="img1" href="http://p.weather.com.cn/2016/11/2625075.shtml"  target="_blank"><img src="http://i.weather.com.cn/images/cn/sjztj/2016/11/27/27095525FBE3A248DDABEEEA04373D4E1448353B.jpg" target="_blank" width="167"  height="126"/><b>广西梧州现雾海奇观 如海市蜃楼</b><i></i></a>

				<a id="img2" href="http://p.weather.com.cn/2016/11/2625019.shtml"  target="_blank"><img src="http://i.weather.com.cn/images/cn/sjztj/2016/11/27/27083557F63946C87C340BFA20877F4C72A71770.jpg" target="_blank" width="167"  height="126"/><b>走过四季的宏村 美如画</b><i></i></a>

				<a id="img3" href="http://p.weather.com.cn/2016/11/2625076.shtml"  target="_blank"><img src="http://i.weather.com.cn/images/cn/sjztj/2016/11/27/271006324526133A9C72B08656D47EDB4998D242.jpg" target="_blank" width="338"  height="255"/><b>重庆阴雨消退迎阳光 市民纷纷晒太阳</b><i></i></a>

				<a id="img4" href="http://p.weather.com.cn/2016/11/2619653.shtml"  target="_blank"><img src="http://pic.weather.com.cn/images/cn/photo/2016/11/17/1710451504BDBF291513C8ABA361CF914B5EB27F.jpg" target="_blank" width="167"  height="126"/><b>盘点全世界那些脑洞大开设计巧妙的大桥</b><i></i></a>

				<a id="img5" href="http://p.weather.com.cn/2016/11/2625107.shtml"  target="_blank"><img src="http://pic.weather.com.cn/images/cn/photo/2016/11/27/2711203165A84CF4DB5E021F2A2E70E4679498D7.jpg" target="_blank" width="167"  height="126"/><b>长春反季光猪节 比基尼美女雪中助阵</b><i></i></a>

			</div>
									<div class="greatEvent">
   <h1>近期重大天气事件</h1>
   <ul>


      <li>
<a href="http://www.weather.com.cn/index/2016/11/2624697.shtml" target="_blank">
	<p class="time">11月26日</p><img src="" width="112" height="84" alt="新疆阿克陶发生6.7级地震 2小时内记录余震57次">
	<div><h2>  新疆阿克陶发生6.7级地震 2小时内记录余震57次  </h2><p>中国地震台网正式测定：11月25日22时24分在新疆克孜勒苏州阿克陶县(北纬39.27度，东经74.04度)发生6.7级地震，震源深度10千米。</p>	</div>
</a>
 </li>


      <li>
<a href="http://www.weather.com.cn/index/2016/11/2624696.shtml" target="_blank">
	<p class="time">11月26日</p><img src="http://i.weather.com.cn/images/cn/news/2016/11/26/260733164220BF4104CEAEA3C0CDE114BAE6426D.jpg" width="112" height="84" alt="台风预警：“蝎虎”最强达台风级 华南沿海迎强风">
	<div><h2>  台风预警：“蝎虎”最强达台风级 华南沿海迎强风  </h2><p>预计，“蝎虎”将以每小时15~20公里的速度向北偏西转偏北方向移动，强度逐渐加强，最强可达强热带风暴级到台风级。27日夜间开始“蝎虎”移速减慢并将在黄岩岛以北的南海中东部海域回旋，28日白天以后将逐渐转向西南方向移动，强度明显减弱乃至消失。</p>	</div>
</a>
 </li>


      <li>
<a href="http://www.weather.com.cn/index/2016/11/2624695.shtml" target="_blank">
	<p class="time">11月26日</p><img src="http://i.weather.com.cn/images/cn/news/2016/11/25/251553345B6AECB2307FAD5001AF435CF0E98864.jpg" width="112" height="84" alt="广东福建局地暴雨 京津冀等地有重度霾">
	<div><h2>  广东福建局地暴雨 京津冀等地有重度霾  </h2><p>京津冀等局地重度霾一张图看哪里霾频发秋冬季注意身体8个变化如何防御突发地质灾害今天（26日），华南、江南南部的强降水将持续，广东、福建局地有暴雨，需防范次生灾害；吉林东部局地有大雪，或对交通产生影响。从今天上午开始，华北、黄淮等地的霾将减弱、消散，28日起还将有雾霾过程。</p>	</div>
</a>
 </li>


      <li>
<a href="http://news.weather.com.cn/2016/11/2623482.shtml" target="_blank">
	<p class="time">11月24日</p><img src="http://i.weather.com.cn/images/cn/news/2016/11/23/23153514D6EA5A3B39C884B2DFFD924BA40AF327.jpg" width="112" height="84" alt="今起我国大部气温回升 华北等地雾霾起">
	<div><h2>  今起我国大部气温回升 华北等地雾霾起  </h2><p>今天（24日），今年下半年来最强寒潮天气过程对我国的影响趋于结束，我国大部地区气温逐渐回升，不过，在冷空气影响结束后，华北、黄淮等地雾霾天气又将发展起来。</p>	</div>
</a>
 </li>


      <li>
<a href="http://news.weather.com.cn/2016/11/2622768.shtml" target="_blank">
	<p class="time">11月23日</p><img src="http://i.weather.com.cn/images/cn/news/2016/11/22/22153934D0FC666AF3377D0CB7C21670CAED4442.jpg" width="112" height="84" alt="湖北安徽等地气温破冰点 两广局地降12℃">
	<div><h2>  湖北安徽等地气温破冰点 两广局地降12℃  </h2><p>今天（23日），冷空气继续影响江南中部及以南地区。预计广西、福建、广东等地的部分地区降温可达12℃以上，湖北、安徽等地局地最低温也将跌破冰点。黄淮南部、江淮和江汉等地仍有降雪。</p>	</div>
</a>
 </li>


      <li>
<a href="http://news.weather.com.cn/2016/11/2622508.shtml" target="_blank">
	<p class="time">11月22日</p><img src="http://i.weather.com.cn/images/cn/news/2016/11/22/22111801E3252B9F726081135906CC6345CD28FB.jpg" width="112" height="84" alt="山西雨雪天气致高速公路多车相撞">
	<div><h2>  山西雨雪天气致高速公路多车相撞  </h2><p>受冷空气影响，昨天（21日），山西部分地区出现雨雪天气，受其影响，山西一高速公路发生多车相撞事故。山西省气象台预计，今天白天，山西南部局部地区仍有小雪，全省气温继续下降，提醒公众注意防范道路结冰对交通的不利影响。</p>	</div>
</a>
 </li>


      <li>
<a href="http://news.weather.com.cn/2016/11/2620715.shtml" target="_blank">
	<p class="time">11月20日</p><img src="http://i.weather.com.cn/images/cn/news/2016/11/19/191533273BB683D761E6E48C9F6FF8A1E2A9BD1F.jpg" width="112" height="84" alt="华北黄淮等迎初雪 中东部大降温若隆冬">
	<div><h2>  华北黄淮等迎初雪 中东部大降温若隆冬  </h2><p>今天（20日）起，一股实力强劲的冷空气开始影响我国，华北、黄淮等地将迎今冬初雪。未来一周，随着冷空气南下，中东部大部地区也将自北向南遭遇显著降温。公众出行需注意雨雪天气的不利影响，并防寒保暖。</p>	</div>
</a>
 </li>


      <li>
<a href="http://news.weather.com.cn/2016/11/2620360.shtml" target="_blank">
	<p class="time">11月19日</p><img src="http://i.weather.com.cn/images/cn/news/2016/11/18/18162511E9634E00E0C2987D6B80C75340E12B85.jpg" width="112" height="84" alt="明起强冷空气速冻中东部 京津冀现初雪">
	<div><h2>  明起强冷空气速冻中东部 京津冀现初雪  </h2><p>预计20-24日，下半年来最强冷空气将影响我国，中东部大部地区将相继出现明显雨雪、大风、降温天气，多地降温16℃以上，气温大面积创下半年来新低，京津冀将迎今冬初雪。</p>	</div>
</a>
 </li>

   </ul>
</div>

			</div>
	<div class="right fr">

						<div class="pic">
   <h1><span><a target="_blank" href="http://p.weather.com.cn" _hover-ignore="1">>></a></span><a href="http://p.weather.com.cn" target="_blank">精彩图集</a></h1>
   <div id='scrollPic' class="m">
      <ul class="bigImg clearfix">


         <li><a href="http://p.weather.com.cn/2016/11/2624864.shtml" target="_blank"><img src="http://pic.weather.com.cn/images/cn/photo/2016/11/26/26141445FACA4050850AF8C60A7D17AE026129DA.jpg" alt="南疆阿合奇雪后现彩虹 惊叹如梦幻" width="300" height="227"></a></li>

         <li><a href="http://p.weather.com.cn/2016/11/2624431.shtml" target="_blank"><img src="http://pic.weather.com.cn/images/cn/photo/2016/11/25/25135856261A45CE96501CEC26E5635A164F18C3.jpg" alt="南方多地上演“冰雪奇缘” 房屋树木披冰甲" width="300" height="227"></a></li>

         <li><a href="http://p.weather.com.cn/2016/11/2624510.shtml" target="_blank"><img src="http://pic.weather.com.cn/images/cn/photo/2016/11/25/251529556987A393ADC4461CFE79D62F0203E6DB.jpg" alt="杭州放晴全民晾晒 “满城尽是被子”" width="300" height="227"></a></li>

         <li><a href="http://p.weather.com.cn/2016/11/2624481.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/public/2016/11/25/251651103209D063E3620C9CA31F2CD8E8693875.jpg" alt="天寒地冻不怕冷的人" width="300" height="227"></a></li>

         <li><a href="http://p.weather.com.cn/2016/08/2580921.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/public/2016/11/23/231018365622699B716EC2C29FCD4E81928D7A51.jpg" alt="-67℃！最寒冷人类居住地的日常生活" width="300" height="227"></a></li>

         <li><a href="http://p.weather.com.cn/2016/11/2617495.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/public/2016/11/24/240901276607F12B57A90AAF8DB3B6905118A373.jpg" alt="设计师智商“欠费”？ 细数那些滑稽设计" width="300" height="227"></a></li>

      </ul>

      <ul class="botIcon clearfix">
      </ul>
      <p></p>
      <div class="bottomBg"></div>
      <div class="rollLeft"></div>
      <div class="rollRight"></div>
   </div>
</div>
<!--焦点图end-->
<ul class="list">
   <li><a href="http://p.weather.com.cn/2016/11/2615872.shtml" target="_blank" style="color:" > 盘点世界上让人叹为观止的自然现象 </a></li>

   <li><a href="http://p.weather.com.cn/2016/11/2623219.shtml" target="_blank" style="color:" > 一场大雪让郑州地铁变成了春运现场 </a></li>

   <li><a href="http://p.weather.com.cn/2016/11/2622374.shtml" target="_blank" style="color:" > 十年雪景！外媒摄影师镜头下的冬日北京 </a></li>

   <li><a href="http://p.weather.com.cn/2016/11/2619202.shtml" target="_blank" style="color:" > 人生若只如初见：与路人40年后“重聚” </a></li>


</ul>


			<div class="ad4" id="adposter_6116">
			<script>
			(function() {
			    var s = "_" + Math.random().toString(36).slice(2);
			    document.write('<div id="' + s + '"></div>');
			    (window.slotbydup=window.slotbydup || []).push({
			        id: '3011990',
			        container: s,
			        size: '300,250',
			        display: 'inlay-fix'
			    });
			})();
			</script>
			</div>

						<div class="chartPH">
				<h1 class="clearfix">
					<span>热点</span>
					<i >视频</i>
					<i>图片</i>
					<i class="on" >文章</i>
				</h1>
				<ul id='hot'>
										<li class="hover"><span class="ord"><i>1</i></span><span class="city"><a  target="_blank" href="http://news.weather.com.cn/2016/11/2619702.shtml">这个冬天格外冷？ 专家预测今冬我国...</a></span></li><li class="hover"><span class="ord"><i>2</i></span><span class="city"><a  target="_blank" href="http://www.weather.com.cn/index/2016/11/2624696.shtml">台风预警：“蝎虎”最强达台风级 华...</a></span></li><li class="hover"><span class="ord"><i>3</i></span><span class="city"><a  target="_blank" href="http://news.weather.com.cn/2016/11/2623482.shtml">今起我国大部气温回升 华北等地雾霾...</a></span></li><li class=""><span class="ord"><i>4</i></span><span class="city"><a  target="_blank" href="http://www.weather.com.cn/life/2016/11/2622757.shtml">入冬防“寒邪” 牢记三字经</a></span></li><li class=""><span class="ord"><i>5</i></span><span class="city"><a  target="_blank" href="http://news.weather.com.cn/2016/11/2622768.shtml">湖北安徽等地气温破冰点 两广局地降...</a></span></li>

									</ul>
				<ul id='pic' style="display:none;">
										<li class="hover"><span class="ord"><i>1</i></span><span class="city"><a  target="_blank" href="http://p.weather.com.cn/2016/08/2580921.shtml">-67℃！最冷村庄的日常生活</a></span></li><li class="hover"><span class="ord"><i>2</i></span><span class="city"><a  target="_blank" href="http://p.weather.com.cn/2016/11/2619653.shtml">盘点全世界那些脑洞大开设计巧妙的大...</a></span></li><li class="hover"><span class="ord"><i>3</i></span><span class="city"><a  target="_blank" href="http://p.weather.com.cn/2016/11/2624481.shtml">天寒地冻不怕冷的人</a></span></li><li class=""><span class="ord"><i>4</i></span><span class="city"><a  target="_blank" href="http://p.weather.com.cn/2016/11/2624510.shtml">杭州放晴全民晾晒 “满城尽是被子”</a></span></li><li class=""><span class="ord"><i>5</i></span><span class="city"><a  target="_blank" href="http://p.weather.com.cn/2016/11/2617495.shtml">当设计师智商“欠费”时 产生了哪些...</a></span></li>

									</ul>
				<ul id="video">
										<li class="hover"><span class="ord"><i>1</i></span><span class="city"><a  target="_blank" href="http://www.weather.com.cn/video/tqyb/05/508815.shtml">新闻联播天气预报</a></span></li><li class="hover"><span class="ord"><i>2</i></span><span class="city"><a  target="_blank" href="http://www.weather.com.cn/video/2016/11/2624565.shtml">这样温度的温泉会泡出危险！</a></span></li><li class="hover"><span class="ord"><i>3</i></span><span class="city"><a  target="_blank" href="http://www.weather.com.cn/video/2016/11/2624593.shtml">今冬会很冷！南方多地降大雪</a></span></li><li class=""><span class="ord"><i>4</i></span><span class="city"><a  target="_blank" href="http://www.weather.com.cn/video/2016/11/2624598.shtml">中央气象台：台风“蝎虎”生成 蓝色...</a></span></li><li class=""><span class="ord"><i>5</i></span><span class="city"><a  target="_blank" href="http://www.weather.com.cn/video/2016/03/lssj/2491187.shtml">百年最强厄尔尼诺将持续至今年5月</a></span></li>

									</ul>
			</div>

			<div class="ad4" id="adposter_6117">
			<script>
(function() {
    var s = "_" + Math.random().toString(36).slice(2);
    document.write('<div id="' + s + '"></div>');
    (window.slotbydup=window.slotbydup || []).push({
        id: '3011991',
        container: s,
        size: '300,250',
        display: 'inlay-fix'
    });
})();
</script>
			</div>


						<div class="pic travel">
				<h1> <span><a href="http://www.weather.com.cn/life/" target="_blank">&gt;&gt;</a></span> <a target="_blank" href="http://www.weather.com.cn/life/">生活旅游</a></h1>
				<div class="scrollPic">
					<ul class="bigImg clearfix">

   <li><a href="http://www.weather.com.cn/life/2016/11/2624432.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/life/2016/11/25/2514053400B5BB43B3DABC9B2D4F7AD56404A49A.jpg" alt="“平安法则”帮老人过冬" width="300" height="227"  ></a></li>

   <li><a href="http://www.weather.com.cn/life/2016/11/2624302.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/life/2016/11/25/251048209432B842AF7FBFB5805FF509C70CFD12.jpg" alt="你比别人老得快的八大原因" width="300" height="227"  ></a></li>

   <li><a href="http://www.weather.com.cn/life/2015/11/gdt/2422588.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/life/2015/11/19/6BB7E89BAFF3C1033A4F6EE38E949D55.jpg" alt="小雪节气雨雪袭 防寒保暖需谨记" width="300" height="227"  ></a></li>



   <li><a href="http://www.weather.com.cn/life/2016/07/2547755.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/life/2016/07/04/04100328FA12175C657B07AE35839F709DBD934D.jpg" alt="那些与海相伴的日子里 你还好吗  " width="300" height="227"></a></li>

   <li><a href="http://www.weather.com.cn/life/2016/06/2540029.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/life/2016/07/04/041007290B2F5A3D889D987A08628FBE778B7058.jpg" alt="在属于骑行的季节里出发" width="300" height="227"></a></li>

   <li><a href="http://www.weather.com.cn/life/2016/05/2522201.shtml" target="_blank"><img src="http://i.weather.com.cn/images/cn/life/2016/05/25/25101117A487BE17A558F9597C5A15969036D84E.jpg" alt="今年夏天去哪里看海？" width="300" height="227"></a></li>


					</ul>

					<ul class="botIcon clearfix">
					</ul>
					<p></p>
					<div class="bottomBg"></div>
					<div class="rollLeft"></div>
					<div class="rollRight"></div>
				</div>
			</div>
			<ul class="list">

   <li><a href="http://www.weather.com.cn/life/2016/05/2521487.shtml" target="_blank">巴黎博物馆“奇妙夜”</a></li>

   <li><a href="http://www.weather.com.cn/life/2016/04/lyxx/2505382.shtml" target="_blank">走进27度的秘密 遇见质朴和自由</a></li>




<li><a href="http://www.weather.com.cn/life/2016/11/2622776.shtml" target="_blank">冬季常熏艾 老少无疾患</a></li>


<li><a href="http://www.weather.com.cn/life/2016/11/2622775.shtml" target="_blank">冬季如何预防冻疮</a></li>

			</ul>

			<div class="ad4" id="adposter_6118">
			<div class="ad4" style="background:url(http://i.tq121.com.cn/i/loading.gif) center center no-repeat #eee"></div>
			</div>



							<div class="hotSpot">
      <h3><b></b>热门景点<i></i></h3>
<h4 class="title"><span class="name">景区</span><span class="weather">天气</span><span class="wd">气温</span><span class="zs">旅游指数</span></h4>
      <ul class="on">
            <li><span class="name"><a href="http://www.weather.com.cn/weather1d/101310201.shtml" target="_blank">三亚</a></span><span class="weather"><a title="多云">多云</a></span><span class="wd">30℃/22℃</span><span class="zs">适宜</span></li><li><span class="name"><a href="http://www.weather.com.cn/weather1d/101271906.shtml" target="_blank">九寨沟</a></span><span class="weather"><a title="多云">多云</a></span><span class="wd">11℃/2℃</span><span class="zs">适宜</span></li><li><span class="name"><a href="http://www.weather.com.cn/weather1d/101290201.shtml" target="_blank">大理</a></span><span class="weather"><a title="晴转多云">晴转多云</a></span><span class="wd">17℃/7℃</span><span class="zs">适宜</span></li><li><span class="name"><a href="http://www.weather.com.cn/weather1d/101251101.shtml" target="_blank">张家界</a></span><span class="weather"><a title="晴">晴</a></span><span class="wd">17℃/4℃</span><span class="zs">适宜</span></li><li><span class="name"><a href="http://www.weather.com.cn/weather1d/101300501.shtml" target="_blank">桂林</a></span><span class="weather"><a title="晴">晴</a></span><span class="wd">18℃/6℃</span><span class="zs">适宜</span></li><li><span class="name"><a href="http://www.weather.com.cn/weather1d/101120201.shtml" target="_blank">青岛</a></span><span class="weather"><a title="晴">晴</a></span><span class="wd">10℃/1℃</span><span class="zs">适宜</span></li>
   </ul>
    </div>




			<div class="ad4" id="adposter_6119">
            <div class="ad4" style="background:url(http://i.tq121.com.cn/i/loading.gif) center center no-repeat #eee"></div>
            </div>




						<div class="weaPro">
			<h1>气象产品</h1>
			<div class="cen clearfix">
			  <div class="l"><a href="http://products.weather.com.cn/product/radar/index/procode/JC_RADAR_AZ9010_JB" target="_blank">
				<img src="http://i.tq121.com.cn/i/weather2014/7d/product.jpg" width='138' height="107" alt=""></a>
				<p><a href="http://products.weather.com.cn/product/radar/index/procode/JC_RADAR_AZ9010_JB" target="_blank">北京基本反射率单站雷达图</a></p>
			  </div>
			  <ul class="r">
				<li><a target="_blank" href="http://www.weather.com.cn/satellite/">中国大陆区域彩色云图</a></li>
				<li><a target="_blank" href="http://www.weather.com.cn/index/zxqxgg1/wlstyb.shtml">未来三天天气趋势预报</a></li>
				<li><a target="_blank" href="http://products.weather.com.cn/product/Index/index/procode/JC_JSL_02405">全国降水量实况</a></li>
				<li><a target="_blank" href="http://products.weather.com.cn/product/Index/index/procode/YB_WD_ZG24">全国最高气温分布</a></li>
				<li><a target="_blank" href="http://products.weather.com.cn/product/Index/index/procode/YB_WD_ZD24">全国最低气温分布</a></li>
			  </ul>
			</div>

			<div class="ad4" id="adposter_6120">
            	<script>
(function() {
    var s = "_" + Math.random().toString(36).slice(2);
    document.write('<div id="' + s + '"></div>');
    (window.slotbydup=window.slotbydup || []).push({
        id: '3011998',
        container: s,
        size: '300,250',
        display: 'inlay-fix'
    });
})();
</script>
            </div>

						<div class="weaSer">
				<h1>气象服务</h1>
				<div>
					<h2>气象服务热线</h2>
					<p>拨打400-6000-121进行气象服务资讯、建议、合作与投诉</p>
				</div>
				<div>
					<h2>天气预报电话查询</h2>
					<p>拨打12121或96121进行天气预报查询</p>
				</div>
				<div>
					<h2>手机查询</h2>
					<p>随时随地通过手机登录中国天气WAP版查看各地天气资讯</p>
				</div>
			</div>

			<div class="ad4" id="adposter_6121">
            <script>
(function() {
    var s = "_" + Math.random().toString(36).slice(2);
    document.write('<div id="' + s + '"></div>');
    (window.slotbydup=window.slotbydup || []).push({
        id: '3012001',
        container: s,
        size: '300,250',
        display: 'inlay-fix'
    });
})();
</script>
            </div>
			</div>
		</div>
</div>
<div id="ab_yjfk">
	<a href="http://www.weather.com.cn/index/feedback_201409.shtml" target="_blank">
		<img border="0" src="http://i.tq121.com.cn/i/weather2014/city/fankui.png" usemap="#Map">
	</a>
</div>
<script>W.use(['j/weather2015/c_common.js','j/weather2015/bluesky/c_1d.js']);</script>


<script>W.js("http://i.tq121.com.cn/j/ad/caoyu-min.js")</script>

<div class="ad-container">


	<div class="adposter_pos" data-posid="adposter_6118">
	<script>
(function() {
    var s = "_" + Math.random().toString(36).slice(2);
    document.write('<div id="' + s + '"></div>');
    (window.slotbydup=window.slotbydup || []).push({
        id: '3011993',
        container: s,
        size: '300,250',
        display: 'inlay-fix'
    });
})();
</script>

	</div>
<div class="adposter_pos" data-posid="adposter_6119">
	<script>
(function() {
    var s = "_" + Math.random().toString(36).slice(2);
    document.write('<div id="' + s + '"></div>');
    (window.slotbydup=window.slotbydup || []).push({
        id: '3011996',
        container: s,
        size: '300,250',
        display: 'inlay-fix'
    });
})();
</script>

	</div>

	<div class="adposter_pos" data-posid="adposter_6298">
	<script>
(function() {
    var s = "_" + Math.random().toString(36).slice(2);
    document.write('<div id="' + s + '"></div>');
    (window.slotbydup=window.slotbydup || []).push({
        id: '3011959',
        container: s,
        size: '680,90',
        display: 'inlay-fix'
    });
})();
</script>

	</div>
	<div class="adposter_pos" data-posid="adposter_6299">
	<script>
(function() {
    var s = "_" + Math.random().toString(36).slice(2);
    document.write('<div id="' + s + '"></div>');
    (window.slotbydup=window.slotbydup || []).push({
        id: '3011963',
        container: s,
        size: '680,90',
        display: 'inlay-fix'
    });
})();
</script>

	</div>
</div>
<!-- START WRating v1.0 -->
	<script type="text/javascript" src="http://m.weather.com.cn/a1.js"></script>
	<script type="text/javascript">
	var vjAcc="860010-2099040100";
	var wrUrl=" http://c.wrating.com/";
	var splitUrl = document.location.toString().split("?");
	if(splitUrl.length == 2){
		if(splitUrl[1].indexOf('from') > -1){
			var search_key = document.location.search.toString().replace("?from=",'');
			if(search_key != 'baidu' && search_key != 'coop' && search_key != 'hao123'){
				search_key = 'weather';
			}
			vjSetReferrer("http://flash."+search_key+".com.cn/come/from/flash.html");
		}
	}
	vjTrack("");
	</script>
	<noscript><img src="http://c.wrating.com/a.gif?a=&c=860010-2099040100" width="1" height="1"/></noscript>
<!-- END WRating v1.0 -->
	<!-- START WRating v1.0 -->
<script type="text/javascript" src="http://click.wrating.com/c3.js"></script>
<script type="text/javascript">
var vjClickAcc="860010-2099330800";
var wrUrl = "http://click.wrating.com/";
initMouseClick();
</script>
<!-- END WRating v1.0 -->
<script>W.use(['j/weather2015/c_gg.js']);</script>
<style>
#abs { z-index:10 ;display:block;}
#adposter_6287{display:block;}
</style>
<script>W.use('j/weather2015/publicHead.js');</script>
<div class="footer">
  <div class="block">
    <div class="Lcontent" style="width:558px;">
      <dl style="width:280px; margin-right:22px;">
        <dt>
          <h3>网站服务</h3>
        </dt>
        <dd>
          <p><a href="http://www.weather.com.cn/wzfw/gywm/">关于我们</a><a href="http://www.weather.com.cn/wzfw/lxwm/">联系我们</a><a href="http://www.weather.com.cn/wzfw/sybz/">帮助</a><a href="http://www.weather.com.cn/wzfw/ryzp/">人员招聘</a></p>
          <p><a href="http://www.weather.com.cn/wzfw/kfzx/">客服中心</a><a href="http://www.weather.com.cn/wzfw/bqsm/">版权声明</a><a href="http://www.weather.com.cn/wzfw/wzls/">律师</a><a href="http://www.weather.com.cn/wzfw/wzdt/">网站地图</a></p>
        </dd>
      </dl>
      <dl style="width:150px;">
        <dt>
          <h3>营销中心</h3>
        </dt>
        <dd>
          <p><a href="http://marketing.weather.com.cn/wzhz/index.shtml">商务合作</a><a href="http://ad.weather.com.cn/index.shtml">广告服务</a></p>
        </dd>
      </dl>
      <div class="clear"></div>
    </div>
    <div class="friendLink" style="width:418px;margin-right:15px;">
      <h3>相关链接</h3>
      <p><a href="http://www.cma.gov.cn/" target="_blank">中国气象局</a><a href="http://pmsc.cma.gov.cn/" target="_blank">公共气象服务中心</a><a href="http://www.chinamsa.org" target="_blank">中国气象服务协会</a> </p>
      <p><a href="http://www.weathertv.cn/" target="_blank">中国气象频道</a><a href="http://www.tourweather.com.cn/" target="_blank">中国旅游天气网</a><a href="http://www.xn121.com/" target="_blank">中国兴农网</a><a target="_blank" href="http://cwera.weather.com.cn/">风能太阳能资源中心</a></p>

    </div>
    <div class="serviceinfo">
<p><span>客服邮箱：<a href="mailto:service@weather.com.cn">service@weather.com.cn</a></span><span style="width:220px;">广告服务：<b>010-58991910</b></span><span><a href="http://www.miibeian.gov.cn/" target="_blank">京ICP证010385号</a>　京公网安备11041400134号</span></p>
      <p><span>客服热线：<b><a href="http://www.weather.com.cn/wzfw/kfzx/index.shtml" target="_blank">400-6000-121</a></b></span><span style="width:220px;">  商务合作：<b>010-58991938</b></span><span>增值电信业务经营许可证B2-20050053</span></p>

    </div>
    <div class="clear"></div>
  </div>
  <div class="aboutUs"> Copyright&copy;<a href="http://pmsc.cma.gov.cn/" target="_blank">中国气象局公共气象服务中心</a> All Rights Reserved (2008-2016) 版权所有 复制必究 郑重声明：中国天气网版权所有，未经书面授权禁止使用 </div>
</div>

<!--<div class="provinceLinks">
  <div class="midBlock" style="margin:0 auto; width:1000px;">
  <dl>
      <dt><a target="_blank" href="http://www.weather.com.cn/textFC/hb.shtml">华北地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://bj.weather.com.cn">北京</a><a target="_blank" href="http://tj.weather.com.cn">天津</a><a target="_blank" href="http://hebei.weather.com.cn">河北</a></p>
        <p><a target="_blank" href="http://shanxi.weather.com.cn">山西</a><a target="_blank" href="http://nmg.weather.com.cn">内蒙古</a></p>
      </dd>
    </dl>
    <dl>
      <dt><a target="_blank" href="http://www.weather.com.cn/textFC/hd.shtml">华东地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://sh.weather.com.cn">上海</a><a target="_blank" href="http://js.weather.com.cn">江苏</a><a target="_blank" href="http://zj.weather.com.cn">浙江</a></p>
        <p><a target="_blank" href="http://ah.weather.com.cn">安徽</a><a target="_blank" href="http://fj.weather.com.cn">福建</a><a target="_blank" href="http://sd.weather.com.cn">山东</a></p>
      </dd>
    </dl>

    <dl>
      <dt><a target="_blank" href="http://www.weather.com.cn/textFC/hz.shtml">华中地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://hubei.weather.com.cn">湖北</a><a target="_blank" href="http://hunan.weather.com.cn">湖南</a></p>
        <p><a target="_blank" href="http://henan.weather.com.cn">河南</a><a target="_blank" href="http://jx.weather.com.cn">江西</a></p>
      </dd>
    </dl>
    <dl>
      <dt><a target="_blank" href="http://www.weather.com.cn/textFC/hn.shtml">华南地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://gd.weather.com.cn">广东</a><a target="_blank" href="http://gx.weather.com.cn">广西</a></p>
        <p><a target="_blank" href="http://hainan.weather.com.cn">海南</a></p>
      </dd>
    </dl>
    <div class="line"></div>
    <dl>
      <dt><a target="_blank" href="http://www.weather.com.cn/textFC/xb.shtml">西北地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://shaanxi.weather.com.cn">陕西</a><a target="_blank" href="http://gs.weather.com.cn">甘肃</a><a target="_blank" href="http://qh.weather.com.cn">青海</a></p>
        <p><a target="_blank" href="http://nx.weather.com.cn">宁夏</a><a target="_blank" href="http://xj.weather.com.cn">新疆</a></p>
      </dd>
    </dl>
    <dl>
      <dt><a target="_blank" href="http://www.weather.com.cn/textFC/xn.shtml">西南地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://cq.weather.com.cn">重庆</a><a target="_blank" href="http://sc.weather.com.cn">四川</a><a target="_blank" href="http://yn.weather.com.cn">云南</a></p>
        <p><a target="_blank" href="http://gz.weather.com.cn">贵州</a><a target="_blank" href="http://xz.weather.com.cn">西藏</a></p>
      </dd>
    </dl>
    <dl>
      <dt><a target="_blank" href="http://www.weather.com.cn/textFC/db.shtml">东北地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://ln.weather.com.cn">辽宁</a><a target="_blank" href="http://jl.weather.com.cn">吉林</a></p>
        <p><a target="_blank" href="http://hlj.weather.com.cn">黑龙江</a></p>
      </dd>
    </dl>
    <dl>
      <dt class="last"><a target="_blank" href="http://www.weather.com.cn/textFC/gat.shtml">港澳台地区</a></dt>
      <dd>
        <p><a target="_blank" href="http://www.weather.com.cn/html/province/xianggang.shtml">香港</a><a target="_blank" href="http://mo.weather.com.cn">澳门</a></p>
        <p><a target="_blank" href="http://www.weather.com.cn/html/province/taiwan.shtml">台湾</a></p>
      </dd>
    </dl>
  </div>


</div>-->
<div style="display:none;">
<script src="http://s11.cnzz.com/z_stat.php?id=1257873453&web_id=1257873453" language="JavaScript"></script>
</div>
<script type="text/javascript">
var _bdhmProtocol = (("https:" == document.location.protocol) ? " https://" : " http://");
document.write(unescape("%3Cscript src='" + _bdhmProtocol + "hm.baidu.com/h.js%3F080dabacb001ad3dc8b9b9049b36d43b' type='text/javascript'%3E%3C/script%3E"));
</script>


<!-- START WRating v1.0 -->
<script type="text/javascript" src="http://c.wrating.com/a1.js">
</script>
<script type="text/javascript">
var vjAcc="860010-1905010101";
var wrUrl="http://c.wrating.com/";
vjTrack("");
</script>
<noscript><img src="http://c.wrating.com/a.gif?a=&c=860010-1905010101" width="1" height="1"/></noscript>
<!-- END WRating v1.0 -->

</body>
</html>`

type A int

func (a A)Write(buf[]byte)(int,error)  {
	fmt.Println(len(buf))
	return len(buf),nil
}

/*
结论是，压缩写非常坑人，在压缩生效时，会将数据拆分成 240 的小块写道底层。
这样配合 tcp NuDelay 会出现底层全是 240 的小包，会出现浪费网络流量的问题。
解决办法有两个，一个是增加在tcp前面再套一个缓冲区，另一个选择是

输出结果：
随机数据：
zlib
2
1
4
16384
1
4
16384
1
4
16384
1
4
16387
1
4
16384
1
4
16384
1
4
4093
1
4
gzip
10
1
4
16384
1
4
16384
1
4
16384
1
4
16387
1
4
16384
1
4
16384
1
4
4093
1
4
deflate
1
4
16384
1
4
16384
1
4
16384
1
4
16387
1
4
16384
1
4
16384
1
4
4093
1
4
zlib:9
2
1
4
16384
1
4
16384
1
4
16384
1
4
16387
1
4
16384
1
4
16384
1
4
4093
1
4
gzip:9
10
1
4
16384
1
4
16384
1
4
16384
1
4
16387
1
4
16384
1
4
16384
1
4
4093
1
4
deflate:9
1
4
16384
1
4
16384
1
4
16384
1
4
16387
1
4
16384
1
4
16384
1
4
4093
1
4
标准html：
zlib
2
240
240
...
240
240
144
4
gzip
10
240
240
240
...
240
144
4
deflate
240
240
240
...
240
240
144
4
zlib:9
2
240
240
240
.。。
240
240
48
4
gzip:9
10
240
...
240
48
4
deflate:9
240
240
...
240
240
240
48
4


*/
func TestSize(t*testing.T){
	zipTypes := []string{"zlib", "gzip", "deflate", "zlib:9", "gzip:9", "deflate:9"}

	f:=func(d []byte){
		for _,zipType:=range zipTypes{
			fmt.Println(zipType)

			a :=A(0)
			zw,err:=NewZipWrite(a,zipType,true)
			if err != nil {
				t.Fatal(err)
			}
			zw.Write(d)
		}
	}

	d:=make([]byte,100*1024)
	rand.Read(d)
	fmt.Println("随机数据：")
	f(d)

	fmt.Println("标准html：")
	d=[]byte(testHtml)
	f(d)
}
/*
随机数据：
zlib
14520
14520
14520
14520
14520
14520
14520
802
gzip
14520
14520
14520
14520
14520
14520
14520
810
deflate
14520
14520
14520
14520
14520
14520
14520
800
zlib:9
14520
14520
14520
14520
14520
14520
14520
802
gzip:9
14520
14520
14520
14520
14520
14520
14520
810
deflate:9
14520
14520
14520
14520
14520
14520
14520
800
标准html：
zlib
14520
5070
gzip
14520
5078
deflate
14520
5068
zlib:9
14520
4974
gzip:9
14520
4982
deflate:9
14520
4972
*/
func TestSize2(t*testing.T){
	zipTypes := []string{"zlib", "gzip", "deflate", "zlib:9", "gzip:9", "deflate:9"}

	f:=func(d []byte){
		for _,zipType:=range zipTypes{
			fmt.Println(zipType)

			a :=A(0)
			b:=bufio.NewWriterSize(a,10*1452)
			zw,err:=NewZipWrite(b,zipType,true)
			if err != nil {
				t.Fatal(err)
			}
			zw.Write(d)
			//b.Flush()
		}
	}

	d:=make([]byte,100*1024)
	rand.Read(d)
	fmt.Println("随机数据：")
	f(d)

	fmt.Println("标准html：")
	d=[]byte(testHtml)
	f(d)
}

func TestNonBlocking(t *testing.T) {
	zipTypes := []string{"zlib", "gzip", "deflate", "zlib:1", "gzip:1", "deflate:1"}

	// 产生随机测试数据
	testDatas := make([][]byte, 1000)
	for i := range testDatas {
		t := make([]byte, rand.Intn(1024)+1)
		rand.Read(t)
		testDatas[i] = t
	}
	l, err := net.Listen("tcp", "127.0.0.1:15634")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	for _, zName := range zipTypes {
		fmt.Println("type:", zName)

		cd := make(chan []byte)
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()

			c, err := l.Accept()
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()

			r, err := NewZipRead(c, zName)
			if err != nil {
				t.Fatal(err)
			}

			for data := range cd {
				buf := make([]byte, len(data))

				if _, err := io.ReadFull(r, buf); err != nil {
					t.Fatal(err)
				} else if bytes.Equal(data, buf) == false {
					t.Fatal(data, "!=", buf)
				}
			}

		}()

		c, err := net.DialTimeout("tcp", l.Addr().String(), 3*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		w, err := NewZipWrite(c, zName,true)
		if err != nil {
			t.Fatal(err)
		}
		for _, data := range testDatas {
			//fmt.Println("write:", len(data))
			if _, err := w.Write(data); err != nil {
				t.Fatal(err)
			}
			//fmt.Println("write ok ")

			cd <- data
		}
		close(cd)

		wg.Wait()
	}
}

func TestBlocking(t *testing.T) {
	zipTypes := []string{"zlib", "gzip", "deflate", "zlib:1", "gzip:1", "deflate:1"}

	// 产生随机测试数据
	testDatas := make([][]byte, 1000)
	for i := range testDatas {
		t := make([]byte, rand.Intn(1024)+1)
		rand.Read(t)
		testDatas[i] = t
	}

	for _, zName := range zipTypes {
		fmt.Println("type:", zName)

		r, w := io.Pipe()

		zr, err := NewZipRead(r, zName)
		if err != nil {
			t.Fatal(err)
		}
		zw, err := NewZipWrite(w, zName,true)
		if err != nil {
			t.Fatal(err)
		}

		c := make(chan []byte)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			for data := range c {
				buf := make([]byte, len(data))
				if _, err := io.ReadFull(zr, buf); err != nil {
					t.Fatal(err)
				} else if bytes.Equal(buf, data) != true {
					t.Fatal(buf, "!=", data)
				}
			}
		}()

		for _, data := range testDatas {
			c <- data
			if _, err := zw.Write(data); err != nil {
				t.Fatal(err)
			}
		}
		close(c)
		wg.Wait()
	}
}
