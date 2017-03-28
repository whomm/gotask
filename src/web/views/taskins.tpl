<!DOCTYPE html>
<html lang="en">
<head>

    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>SB Admin 2 - Bootstrap Admin Theme</title>

    <!-- Bootstrap Core CSS -->
    <link href="/static/vendor/bootstrap/css/bootstrap.min.css" rel="stylesheet">

    <!-- MetisMenu CSS -->
    <link href="/static/vendor/metisMenu/metisMenu.min.css" rel="stylesheet">

    <!-- Custom CSS -->
    <link href="/static/dist/css/sb-admin-2.css" rel="stylesheet">

    <!-- Custom Fonts -->
    <link href="/static/vendor/font-awesome/css/font-awesome.min.css" rel="stylesheet" type="text/css">

    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
        <script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->

</head>

<body>
    <div id="wrapper">


       
            <div class="container-fluid">
                <div class="row">
                    <div class="col-lg-12">
                        <h1 class="page-header"></h1>
                        <div class="panel panel-default">
                        <div class="panel-heading">
                            {{.taskinfo.Name}}
                        </div>
                        <!-- /.panel-heading -->
                        <div class="panel-body">

                            <div class="row form-inline" >
                                <div class="col-sm-6" >
                                    <div class="dataTables_length" >
                                        <form action="" method="get">
                                            starttime: endtime:
                                            <button type="button" class="btn btn-primary btn-sm" onclick="window.location.href='/taskupdate';">Create</button>
                                        </form>
                                    </div>
                                </div>
                                <div class="col-sm-6" style="text-align: right">
                                    
                                </div>
                            </div>

                            <div class="table-responsive">
                                <table class="table  table-striped table-bordered table-hover dataTable no-footer dtr-inline">
                                    <thead>
                                        <tr>
                                            <th>#</th>
                                            <th>Runtime</th>
                                            <th>Tasktime</th>
                                            <th>Status</th>
                                            <th>Createby</th>
                                            <th></th>
                                            <td>Calltime</td>
                                          
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {{range .list}}
                                        <tr>
                                            <td>{{.Id}}</td>
                                            <td>{{ .Runtime | int64todate }}</td>
                                            <td>{{ .Tasktime | int64todate }}</td>
                                            <td>{{if eq .Status  0 }}     PENGDING 
                                                {{else if eq .Status  1}} READY
                                                {{else if eq .Status  2}} RUN
                                                {{else if eq .Status  3}} SUCCESS
                                                {{else if eq .Status  4}} FAIL
                                                {{else if eq .Status  5}} KILLING
                                                {{else if eq .Status  6}} KILLED
                                                {{else if eq .Status  7}} CALLING
                                                {{else if eq .Status  8}} CALLFAIL
                                                {{else if eq .Status  9}} READY
                                                {{end}}
                                            </td>
                                            <td>{{if le .Createby  0 }} sys {{else}} ww {{end}}</td>   
                                            <td>
                                               {{if gt .Time_create  0 }} create: {{ .Time_create|int64todate }} / {{end}} 
                                              {{if gt .Time_ready  0 }}ready: {{ .Time_ready|int64todate }}/{{end}}
                                              {{if gt .Time_run  0 }}run: {{ .Time_run|int64todate }}/{{end}}
                                              {{if gt .Time_success  0 }} success: {{ .Time_success|int64todate }}/{{end}}
                                              {{if gt .Time_fail  0 }}fail: {{ .Time_fail|int64todate }}/{{end}}
                                              {{if gt .Time_callfail  0 }}callfail: {{ .Time_callfail|int64todate }}/{{end}}
                                              {{if gt .Time_killing  0 }}killing: {{ .Time_killing|int64todate }}/{{end}}
                                              {{if gt .Time_killed  0 }}killed: {{ .Time_killed|int64todate }}{{end}}
                                              </td>
                                            <td>{{.Calltime}}</td>
                                        </tr>
                                        {{end}}
                                        
                                    </tbody>
                                    
                                        
                                      
                                </table>
                                
                            </div>
                            <!-- /.table-responsive -->
                            <div class="row">
                                <div class="col-lg-6">
                                    <div class="dataTables_info"  role="status" aria-live="polite">
                                    {{ .pageinfo}}
                                    </div>
                                </div>
                                <div class="col-lg-6">
                                    <div  style="white-space: nowrap; text-align: right;">
                                        <ul class="pagination" style="margin: 0;">
                                            <li class="paginate_button previous">
                                                <a href="{{.pagearr.Begin}}">Begin</a>
                                            </li>
                                            <li class="paginate_button" >
                                                <a href="{{.pagearr.Previous}}">Previous</a>
                                            </li>
                                            <li class="paginate_button">
                                                <a href="{{.pagearr.Next}}">Next</a>
                                            </li>
                                            <li class="paginate_button next">
                                                <a href="{{.pagearr.End}}">End</a>
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                        
                            
                     
                    </div>
                    </div>
                    <!-- /.col-lg-12 -->
                </div>
                <!-- /.row -->
            </div>
            <!-- /.container-fluid -->
       

    </div>
    <!-- /#wrapper -->

    <!-- jQuery -->
    <script src="/static/vendor/jquery/jquery.min.js"></script>

    <!-- Bootstrap Core JavaScript -->
    <script src="/static/vendor/bootstrap/js/bootstrap.min.js"></script>

</body></html>