# TOKO IJAH

Toko Ijah is a web application developed using Go.
It manages the inventory of Toko Ijah.

# Getting Started

We can use $PORT variable to set the page port.
Default port is 9876 (http://localhost:9876).

``` console
go get
go run main.go
```

# Tables/Listings

Toko Ijah consists of three main tables/listings:
* Catatan Jumlah Barang (ItemAmount)
* Catatan Barang Masuk (ItemIn)
* Catatan Barang Keluar (ItemOut)

Each tables can be modified using CRUD endpoints and can be imported using csv file.
Both ItemIn and ItemOut tables are related to ItemAmount linked to SKU.

## Catatan Jumlah Barang (ItemAmount)

Name | Form Data Name | Type
-------- | ------------- | -------
SKU | SKU | string
Nama Item | Name | string
Jumlah Sekarang | Quantity | integer

## Catatan Barang Masuk (ItemIn)

Name | Form Data Name | Type
-------- | ------------- | -------
ID | ID | integer
Waktu | Time | string
SKU | SKU | string
Nama Barang | Name | string
Jumlah Pemesanan | AmountOrders | integer
Jumlah Diterima | AmountReceived | integer
Harga Beli | PurchasePrice | integer
Total | (N/A) | integer
Nomer Kwitansi | ReceiptNumber | string
Catatan | Notes | string

The export result for this table has an additional ID column to identify the rows easier.
But to import the table from csv, it uses the structure of the original one.

## Catatan Barang Keluar (ItemOut)

Name | Form Data Name | Type
-------- | ------------- | -------
ID | ID | integer
Waktu | Time | string
SKU | SKU | string
Nama Barang | Name | string
Jumlah Keluar | AmountOut | integer
Harga Jual | SellingPrice | integer
Total | (N/A) | integer
ID Pesanan | OrderID | string
Catatan | Notes | string

This table has a difference with the original one: the ID Pesana (Order ID) column.
It is created because the value is needed to generate Laporan Penjualan.

The export result for this table has an additional ID column to identify the rows easier, and also the Order ID.
But to import the table from csv, it uses the structure of the original one.

# Reports

Toko Ijah can export two type of reports:
* Laporan Nilai Barang (ItemValueReport)
* Laporan Penjualan (SellingReport)
Each reports can be exported as csv files.

To generate selling report, we need to input some data form for filtering:

Name | Form Data Name | Type
-------- | ------------- | -------
Tanggal Mulai | DateFrom | integer (dd MMMM yyyy format)
Tanggal Akhir | DateTo | string (dd MMMM yyyy format)

NOTES: "Jumlah" column in Laporan Nilai Barang is generated as the sum of items from Catatan Barang Masuk

# Web Page

Toko Ijah has a homepage that is linked to those 5 listings:
* Catatan Jumlah Barang
* Catatan Barang Masuk
* Catatan Barang Keluar
* Laporan Nilai Barang
* Laporan Penjualan

Every pages have features to fullfil the RESTFUL endpoints.

# Other Assumptions

To make things easier, it is assumed that:
* SKU of items won't be changed so SKU is used as primary key
* SKU is generated either manually by the user or by another tool
* Amount of Item, amount of selling, amount of purchase will be updated manually by the user so this web app will not handle the item amount
* CSV files that are used to import data are already valid (the same as the existing spreadsheet)