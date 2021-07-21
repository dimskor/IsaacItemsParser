<?php
$json = file_get_contents('items.json');
$data = json_decode($json, true);
?>

<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Items</title>
	<style>
		body {
			font-family: system-ui,-apple-system,"Segoe UI",Roboto,"Helvetica Neue",Arial,"Noto Sans","Liberation Sans",sans-serif;
		}
		table {
			border-collapse: collapse;
		}
		th, td {
			border: solid 1px grey;
			padding: 4px;
		}
		img {
			max-width: 45px;
			vertical-align: top;
		}
	</style>
</head>
<body>
	<table>
		<tr>
			<th>Id</th>
			<th>Icon</th>
			<th>Name</th>
			<th>Charges</th>
			<th>Quote</th>
			<th>Description</th>
			<th>Quality</th>
			<th>Type</th>
		</tr>
		<?php foreach ($data as $item): ?>
			<tr>
				<td><?php echo $item['Id']; ?></td>
				<td><img src="<?php echo $item['IconUrl']; ?>"></td>
				<td>
					<a href="<?php echo $item['PageUrl']; ?>"><?php echo $item['Name']; ?></a>
				</td>
				<td><?php echo $item['Charges'] ?? '-'; ?></td>
				<td><?php echo $item['Quote']; ?></td>
				<td><?php echo nl2br($item['Description']); ?></td>
				<td><?php echo $item['Quality'] ?? '?'; ?></td>
				<td><?php echo $item['Type']; ?></td>
			</tr>
		<?php endforeach ?>
	</table>
</body>
</html>
