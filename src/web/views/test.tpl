
<!DOCTYPE html>
<html>
<head>
  <title>Beego</title>
    <!-- 最新版本的 Bootstrap 核心 CSS 文件 -->
    <link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
    <!-- 可选的 Bootstrap 主题文件（一般不用引入） -->
    <link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">
</head>
<body>


<div class="container">
  <div class="row">
    <div class="col-xs-3" id="mylist">
    </div>
    <div class="col-xs-9" id="mycontent" >
     
    </div>
  </div>
</div>

</body>
<!-- 最新的 Bootstrap 核心 JavaScript 文件 -->
<script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
<script src="//cdn.bootcss.com/jquery/1.6.1/jquery.js"></script>
<script>
var myforms=
[

   {
        title:"任务组创建/更新",
        method:"get",
        action:"/v1api/tgupdate",
        input:
        [
            {name:"id",value:"0",des:"必须 0创建其他更新  uint64"},
            {name:"uid",value:"1",des:"必须 所有者  uint64"},
            {name:"name",value:"测试任务",des:"必须 任务名称 varchar(60) "},
            {name:"extra",value:"{}",des:"必须   配置json格式"}

        ],
        info:""
    },
    {
        title:"任务组删除",
        method:"get",
        action:"/v1api/tgremove",
        input:
        [
            {name:"id",value:"1",des:"必须    uint64"}


        ],
        info:"组内无任务可以删除"
    },
    {
        title:"任务组查询",
        method:"get",
        action:"/v1api/tgsearch",
        input:
        [
            {name:"id", value:"", des:"可选 不能为空 任务组id"},
            {name:"name",value:"",des:"可选 不能为空 任务名称"},
            /*{name:"ugid",value:"",des:"可选 用户组id"},*/
            {name:"uid",value:"1",des:"可选 不能为空 用户id"},
            {name:"ps",value:"0",des:"必须 开始页面"},
            {name:"pl",value:"10",des:"必须 页面大小"}

        ],
        info:""
    },

    {
        title:"任务创建/更新",
        method:"get",
        action:"/v1api/tupdate",
        input:
        [
            {name:"id",value:"0",des:"必须 0创建其他更新  uint64"},
            {name:"uid",value:"1",des:"必须 所有者  uint64"},
            {name:"tgid",value:"1",des:"必须 任务组id  uint64"},
            {name:"name",value:"测试任务",des:"必须 任务名称 varchar(60) "},
            {name:"crontab",value:"5 * * * * *",des:"必须  调度设置（空为不定期运行）"},
            {name:"pendingtime",value:"30",des:"必须  延迟启动时间/秒"},
            {name:"starttime",value:"0",des:"必须  时间戳 开始时间"},
            {name:"endtime",value:"1483334710",des:"必须 时间戳 结束时间"},
            {name:"extra",value:"{\"workrpc\":\"http://127.0.0.1:8912/\"}",des:"必须  配置json格式"},
            {name:"invalid",value:"0",des:"必须 [0,1]  是否无效"},
            {name:"relay",value:"{\"rl\":[]}",des:"必须 依赖任务列表{id: 3434 start: -1day length: 1day}"}


        ],
        info:""
    },
    {
        title:"任务删除",
        method:"get",
        action:"/v1api/tremove",
        input:
        [
            {name:"id",value:"10",des:"必须  任务id"}
        ],
        info:""
    },

    {
        title:"任务查询",
        method:"get",
        action:"/v1api/tsearch",
        input:
        [
            {name:"id",value:"13121905687",des:"可选 不能为空 任务id"},
            {name:"name",value:"",des:"可选 不能为空 任务名称"},
            {name:"tgid",value:"",des:"可选 不能为空  任务组id"},
            /*{name:"ugid",value:"",des:"必须 可以为空  用户组id"},*/
            {name:"uid",value:"",des:"可选  不能为空 用户id"},
            {name:"ps",value:"0",des:"必须 开始页面"},
            {name:"pl",value:"10",des:"必须 页面大小"}

        ],
        info:""
    },
    {
        title:"任务运行记录",
        method:"get",
        action:"/v1api/thistory",
        input:
        [
            {name:"id",value:"13121905687",des:"必须 任务id"},
            {name:"ps",value:"0",des:"必须 开始页面"},
            {name:"pl",value:"10",des:"必须 页面大小"}
        ],
        info:""
    },
    {
        title:"启动任务",
        method:"get",
        action:"/v1api/tinstance",
        input:
        [
            {name:"id",value:"13121905687",des:"必须 任务id"},
            {name:"starttime",value:"13121905687",des:"必须 时间戳"},
            {name:"endtime",value:"13121905687",des:"必须 时间戳"}
        ],
        info:"一个关联时间同时只能有一个任务实例处于活动状态,返回启动的任务数目"
    },
    {
        title:"杀死任务",
        method:"get",
        action:"/v1api/tkill",
        input:
        [
            {name:"id",value:"13121905687",des:"必须 任务实例id"}
        ],
        info:""
    },


    {
        title:"[cb]任务运行完成",
        method:"get",
        action:"/v1api/insstatsfish",
        input:
        [
            {name:"id",value:"13121905687",des:"必须 任务实例id"}
        ],
        info:""
    },
    {
        title:"[cb]任务运行失败",
        method:"get",
        action:"/v1api/insstatsfal",
        input:
        [
            {name:"id",value:"13121905687",des:"必须 任务实例id"}
        ],
        info:""
    },

    {
        title:"[cb]任务杀死",
        method:"get",
        action:"/v1api/insstatskilled",
        input:
        [
            {name:"id",value:"13121905687",des:"必须 任务实例id"}
        ],
        info:""
    }
];



$(document).ready(function(){


        var html='';
        var listhtml='';
          $.each(myforms,function(i,n){

             var enc = '';
             if(n.enctype)
             {
                enc = "enctype="+n.enctype;
             }
             listhtml +=
             "<li><a href=\"#"+n.title+"\">"+i+'.'+n.title+"</a></li>";

             html +=
             "<div>"+
             "<a name=\""+n.title+"\" id=\""+n.title+"\" ></a>"+
             "<h3>"+i+"."+n.title+"</h3><hr>"+
                "<p>调用接口："+n.action+"</p>"+
                "<p>返回说明："+n.info+"</p>"+
                "<form target=\"_blank\" class=\"well form-inline\" id="+n.title+" "+enc+" action="+n.action+" method="+n.method+">"+
                     createinput(n.input)+
                    '<input type="submit" class="btn btn-default">'+
                "</form>"+
             "<div>";

         });

         $('#mycontent').html(html);
         $('#mylist').html(listhtml);
         function createinput(obj)
         {
            if(obj.length>0)
            {
                var ret='';
                $.each(obj,function(i,n){

                    var des = '';
                    if(n.des)
                    {
                        des = n.des;
                    }
                    if(n.type)
                    {
                        ret += '<li><span style=\"display:inline-block;width:100px;\">'+n.name+':</span><input type="'+n.type+'" name="'+n.name+'" value="'+n.value+'" />'+des+'</li>';
                    }
                    else
                    {
                        ret += '<li><span style=\"display:inline-block;width:100px;\">'+n.name+':</span><input type="text" name="'+n.name+'" value=\''+n.value+'\'/>'+des+'</li>';
                    }
                });
                return ret;
            }
            return '';
         }



});
</script>
</html>
