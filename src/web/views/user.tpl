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
                            User list
                        </div>
                        <!-- /.panel-heading -->
                        <div class="panel-body">

                            <div class="row form-inline" >
                                <div class="col-sm-6" >
                                    <div class="dataTables_length" >
                                        <button type="button" class="btn btn-primary btn-sm">Create</button>
                                        <!--
                                        <label>Show 
                                            <select name="dataTables-example_length"  class="form-control input-sm">
                                                <option value="10">10</option>
                                                <option value="25">25</option>
                                                <option value="50">50</option>
                                                <option value="100">100</option>
                                            </select> entries
                                        </label>
                                        -->
                                    </div>
                                </div>
                                <div class="col-sm-6" style="text-align: right">
                                    <form action="" method="get">
                                    <div id="dataTables-example_filter" class="dataTables_filter">
                                        <label>Search:
                                            <input type="search" class="form-control input-sm" placeholder="" name="keyword" >
                                        </label>
                                    </div>
                                    </form>
                                </div>
                            </div>

                            <div class="table-responsive">
                                <table class="table  table-striped table-bordered table-hover dataTable no-footer dtr-inline">
                                    <thead>
                                        <tr>
                                            <th>#</th>
                                            <th>Name</th>
                                            <th>Usergourp</th>
                                            <th>Edit</th>
                                 
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {{range .list}}
                                        <tr>
                                            <td>{{.Id}}</td>
                                            <td>{{.Name}}</td>
                                            <td>{{.Ugid}}</td>
                                            <td><a href="/#?id={{.Id}}">Edit</a></td>
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