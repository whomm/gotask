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
    <link href="/static/vendor/bootstrap/css/bootstrap.css" rel="stylesheet">

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

    <link rel="stylesheet" href="/static/css/bootstrap-datetimepicker.min.css" />

</head>

<body>

          
    <div id="wrapper">
            <div class="container-fluid">
              <div class="row">
                <div class="col-lg-12">
                    <h1 class="page-header"></h1>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->
            <div class="row">
                <div class="col-lg-12">
                    <div class="panel panel-default">
                        <div class="panel-heading">
                            Create a task
                        </div>
                        <div class="panel-body">
                             <div class="row">
                                <div class="col-lg-12">

                                
                                     <form role="form" action="/tasksave" method="get" >
                                      <div class="form-group">
                                            <label>Name:</label>
                                            <input class="form-control" placeholder="OnlineUserCount" name="name">
                                            <p class="help-block">Task name.</p>
                                        </div>
                                        <div class="form-group">
                                            <label>Group:</label>
                                            
                                            <select class="form-control" name="tgid">
                                            {{range .tglist}}
                                                <option value="{{.Id}}">{{.Name}}</option>
                                            {{end}}
                                            </select>
                                            <p class="help-block">Chose a group to classify tasks.</p>
                                        </div>
                                        <div class="form-group">
                                            <label>Crontab</label>
                                            <input class="form-control" placeholder="5 * * * * *" name="crontab">
                                            <p class="help-block">Same as crontab.</p>
                                        </div>
                                        <div class="form-group">
                                            <label>Pending</label>
                                            <input class="form-control" placeholder="30" name="pendingtime">
                                            <p class="help-block">wait seconds after crontab set</p>
                                        </div>


                                       <div class="form-group">
                                            <label>Start</label>
                                            <div class='input-group date' id='starttime'>
                                                <input type='text' class="form-control" name="starttime" />
                                                <span class="input-group-addon">
                                                    <span class="glyphicon glyphicon-calendar"></span>
                                                </span>
                                            </div>
                                        </div>
                                        <div class="form-group">
                                            <label>End</label>
                                            <div class='input-group date' id='endtime'>
                                                <input type='text' class="form-control" name="endtime" />
                                                <span class="input-group-addon">
                                                    <span class="glyphicon glyphicon-calendar"></span>
                                                </span>
                                            </div>
                                        </div>
                                   
                                        

                                        
                                        <div class="form-group extform">
                                            <label>Extra</label>
                                            {
                                            <div class="extitem form-group input-group">
                                                <div class=" input-group">
                                                    <input class="form-control" name="extkey[]" type="text" placeholder="workrpc" />
                                                    <span class="input-group-addon">:</span>
                                                    <input class="form-control" name="extval[]" type="text" placeholder="http://127.0.0.1:8912" />
                                                    <span class="input-group-btn">
                                                        <button class="btn btn-default btn-add" type="button">
                                                            <span class="glyphicon glyphicon-plus"></span>
                                                        </button>
                                                    </span>
                                                </div>
                                                
                                            </div>
                                          
                                            
                                        </div>
                                        }
                                        

                                        <div class="form-group">
                                            <label>Invalid</label>
                                            <label class="radio-inline">
                                                <input type="radio" name="invalid" id="optionsRadiosInline1" value="0" checked="">False
                                            </label>
                                            <label class="radio-inline">
                                                <input type="radio" name="invalid" id="optionsRadiosInline2" value="1">True
                                            </label>
                                           
                                        </div>

                                        <div class="form-group">
                                            <label>Relay</label>
                                            <textarea class="form-control" rows="3" name="relay" value='{"rl":[]}'></textarea>
                                        </div>

                                        
                                        <button type="reset" class="btn btn-default">Reset</button>
                                        <button type="submit" class="btn btn-default">Save</button>
                                    </form>
                    
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



  
  <script type="text/javascript" src="/static/js/moment.min.js"></script>
  <script type="text/javascript" src="/static/js/bootstrap-datetimepicker.min.js"></script>
  
    <script type="text/javascript">
        $(function () {



            $('#starttime').datetimepicker({format: 'YYYY-MM-DD HH:mm:ss'});
            $('#endtime').datetimepicker({
                format: 'YYYY-MM-DD HH:mm:ss',
                useCurrent: false //Important! See issue #1075
            });
            $("#starttime").on("dp.change", function (e) {
                $('#endtime').data("DateTimePicker").minDate(e.date);
            });
            $("#endtime").on("dp.change", function (e) {
                $('#starttime').data("DateTimePicker").maxDate(e.date);
            });




            $(document).on('click', '.btn-add', function(e)
            {
                e.preventDefault();
               
                var controlForm = $('.extform '),
                    currentEntry = $(this).parents().parents('.extitem:first'),
                    newEntry = $(currentEntry.clone()).appendTo(controlForm);

                newEntry.find('input').val('');
                controlForm.find('.extitem:not(:last) .btn-add')
                    .removeClass('btn-add').addClass('btn-remove').html('<span class="glyphicon glyphicon-minus"></span>');
                   
                    
            }).on('click', '.btn-remove', function(e)
            {
                $(this).parents('.extitem:first').remove();

                e.preventDefault();
                return false;
            });

        });
    </script>
</body></html>