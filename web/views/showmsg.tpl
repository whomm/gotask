<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<title>提示信息</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
<div id="main">
<h2>提示信息</h2>
<table class="showmessage">
	<tbody>
		<tr>
			<td height="30"><p>{{.msg}}</p></td>
		</tr>
	</tbody>
	<tfoot>
	<tr>
		<td colspan="20" align="center">
		<a href="{{.redirect}}">3秒后自动返回上一页面...</a>
		<script type="text/javascript">
			setTimeout("window.location.href='{{.redirect}}'", 3000);
		</script>
		</td>
	</tr>
	</tfoot>
</table>
</div>
</body>
</html>