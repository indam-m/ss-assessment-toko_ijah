# TOKO IJAH

Toko Ijah is a web application developed using Go.

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

# Reports

Toko Ijah can export two type of reports:
* Laporan Nilai Barang (ItemValueReport)
* Laporan Penjualan (SellingReport)
Each reports can be exported as csv files.

To generate selling report, we need to input some data form:

Name | Form Data Name | Type
-------- | ------------- | -------
Tanggal Mulai | DateFrom | integer (dd MMMM yyyy format)
Tanggal Akhir | DateTo | string (dd MMMM yyyy format)