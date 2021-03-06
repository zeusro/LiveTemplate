
1. [追加文本内容](#追加文本内容)
1. [NPOI读写excel](#NPOI读写excel)
1. [Aspose.Cells操作excel](#Aspose.Cells操作excel)
1. [防止死锁的运行异步方法](#防止死锁的运行异步方法)
1. [导出C#生成的公私钥](#导出C#生成的公私钥)
1. [将时间字符串转为Json时间](#将时间字符串转为Json时间)
1. [](#)

## 追加文本内容
```csharp
  using (var fs = File.OpenWrite(Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "log", DateTime.Now.ToString("yyyy-MM-dd") + ".log")))
            {
                //设定书写的开始位置为文件的末尾  
                fs.Position = fs.Length;
                //将待写入内容追加到文件末尾  
                fs.Write(bytes, 0, bytes.Length);
            }
```

## NPOI读写excel
```csharp
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using NPOI.SS.UserModel;
using NPOI.XSSF.UserModel;
using NPOI.HSSF.UserModel;
using System.IO;
using System.Data;

namespace NetUtilityLib
{
    public class ExcelHelper : IDisposable
    {
        private string fileName = null; //文件名
        private IWorkbook workbook = null;
        private FileStream fs = null;
        private bool disposed;

        public ExcelHelper(string fileName)
        {
            this.fileName = fileName;
            disposed = false;
        }

        /// <summary>
        /// 将DataTable数据导入到excel中
        /// </summary>
        /// <param name="data">要导入的数据</param>
        /// <param name="isColumnWritten">DataTable的列名是否要导入</param>
        /// <param name="sheetName">要导入的excel的sheet的名称</param>
        /// <returns>导入数据行数(包含列名那一行)</returns>
        public int DataTableToExcel(DataTable data, string sheetName, bool isColumnWritten)
        {
            int i = 0;
            int j = 0;
            int count = 0;
            ISheet sheet = null;

            fs = new FileStream(fileName, FileMode.OpenOrCreate, FileAccess.ReadWrite);
            if (fileName.IndexOf(".xlsx") > 0) // 2007版本
                workbook = new XSSFWorkbook();
            else if (fileName.IndexOf(".xls") > 0) // 2003版本
                workbook = new HSSFWorkbook();

            try
            {
                if (workbook != null)
                {
                    sheet = workbook.CreateSheet(sheetName);
                }
                else
                {
                    return -1;
                }

                if (isColumnWritten == true) //写入DataTable的列名
                {
                    IRow row = sheet.CreateRow(0);
                    for (j = 0; j < data.Columns.Count; ++j)
                    {
                        row.CreateCell(j).SetCellValue(data.Columns[j].ColumnName);
                    }
                    count = 1;
                }
                else
                {
                    count = 0;
                }

                for (i = 0; i < data.Rows.Count; ++i)
                {
                    IRow row = sheet.CreateRow(count);
                    for (j = 0; j < data.Columns.Count; ++j)
                    {
                        row.CreateCell(j).SetCellValue(data.Rows[i][j].ToString());
                    }
                    ++count;
                }
                workbook.Write(fs); //写入到excel
                return count;
            }
            catch (Exception ex)
            {
                Console.WriteLine("Exception: " + ex.Message);
                return -1;
            }
        }

        /// <summary>
        /// 将excel中的数据导入到DataTable中
        /// </summary>
        /// <param name="sheetName">excel工作薄sheet的名称</param>
        /// <param name="isFirstRowColumn">第一行是否是DataTable的列名</param>
        /// <returns>返回的DataTable</returns>
        public DataTable ExcelToDataTable(string sheetName, bool isFirstRowColumn)
        {
            ISheet sheet = null;
            DataTable data = new DataTable();
            int startRow = 0;
            try
            {
                fs = new FileStream(fileName, FileMode.Open, FileAccess.Read);
                if (fileName.IndexOf(".xlsx") > 0) // 2007版本
                    workbook = new XSSFWorkbook(fs);
                else if (fileName.IndexOf(".xls") > 0) // 2003版本
                    workbook = new HSSFWorkbook(fs);

                if (sheetName != null)
                {
                    sheet = workbook.GetSheet(sheetName);
                    if (sheet == null) //如果没有找到指定的sheetName对应的sheet，则尝试获取第一个sheet
                    {
                        sheet = workbook.GetSheetAt(0);
                    }
                }
                else
                {
                    sheet = workbook.GetSheetAt(0);
                }
                if (sheet != null)
                {
                    IRow firstRow = sheet.GetRow(0);
                    int cellCount = firstRow.LastCellNum; //一行最后一个cell的编号 即总的列数

                    if (isFirstRowColumn)
                    {
                        for (int i = firstRow.FirstCellNum; i < cellCount; ++i)
                        {
                            ICell cell = firstRow.GetCell(i);
                            if (cell != null)
                            {
                                string cellValue = cell.StringCellValue;
                                if (cellValue != null)
                                {
                                    DataColumn column = new DataColumn(cellValue);
                                    data.Columns.Add(column);
                                }
                            }
                        }
                        startRow = sheet.FirstRowNum + 1;
                    }
                    else
                    {
                        startRow = sheet.FirstRowNum;
                    }

                    //最后一列的标号
                    int rowCount = sheet.LastRowNum;
                    for (int i = startRow; i <= rowCount; ++i)
                    {
                        IRow row = sheet.GetRow(i);
                        if (row == null) continue; //没有数据的行默认是null　　　　　　　
                        
                        DataRow dataRow = data.NewRow();
                        for (int j = row.FirstCellNum; j < cellCount; ++j)
                        {
                            if (row.GetCell(j) != null) //同理，没有数据的单元格都默认是null
                                dataRow[j] = row.GetCell(j).ToString();
                        }
                        data.Rows.Add(dataRow);
                    }
                }

                return data;
            }
            catch (Exception ex)
            {
                Console.WriteLine("Exception: " + ex.Message);
                return null;
            }
        }

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (!this.disposed)
            {
                if (disposing)
                {
                    if (fs != null)
                        fs.Close();
                }

                fs = null;
                disposed = true;
            }
        }
    }
}
```

## Aspose.Cells操作excel
```csharp
using System;
using System.Collections.Generic;
using System.Data;
using System.IO;
using System.Linq;
using System.Text;
using Aspose.Cells;

namespace NetUtilityLib
{
    public static class ExcelHelper
    {
        public static int DataTableToExcel(DataTable data, string fileName, string sheetName, bool isColumnNameWritten)
        {
            int num = -1;
            try
            {
                Workbook workBook;
                Worksheet worksheet = null;

                if (File.Exists(fileName))
                    workBook = new Workbook(fileName);
                else
                    workBook = new Workbook();

                if (sheetName == null)
                {
                    if (workBook.Worksheets.Count > 0)
                    {
                        worksheet = workBook.Worksheets[0];
                    }
                    else
                    {
                        sheetName = "Sheet1";
                        workBook.Worksheets.RemoveAt(sheetName);
                        worksheet = workBook.Worksheets.Add(sheetName);
                    }
                }
                if (worksheet != null)
                {
                    worksheet.Cells.Clear();
                    num = worksheet.Cells.ImportDataTable(data, isColumnNameWritten, 0, 0, false);
                    workBook.Save(fileName);
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }

            return num;
        }

        public static void AddOneRowToExcel(DataRow dataRow, string fileName, string sheetName)
        {
            try
            {
                Workbook workBook;

                if (File.Exists(fileName))
                    workBook = new Workbook(fileName);
                else
                    workBook = new Workbook();

                Worksheet worksheet=null;

                if (sheetName == null)
                {
                    worksheet = workBook.Worksheets[0];
                }
                else
                {
                    worksheet = workBook.Worksheets[sheetName];
                }
                if (worksheet != null)
                {
                    worksheet.Cells.ImportDataRow(dataRow, worksheet.Cells.MaxDataRow + 1,0);
                    //worksheet.Cells.ImportArray(dataArray, worksheet.Cells.MaxDataRow+1, 0, false);
                    workBook.Save(fileName);
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

        public static DataTable ExcelToDataTable(string fileName, string sheetName, bool isFirstRowColumnName)
        {
            DataTable data = new DataTable();

            try
            {
                Workbook workbook = null;

                FileInfo fileInfo = new FileInfo(fileName);
                if (fileInfo.Extension.ToLower().Equals(".xlsx"))
                    workbook = new Workbook(fileName, new LoadOptions(LoadFormat.Xlsx));
                else if (fileInfo.Extension.ToLower().Equals(".xls"))
                    workbook = new Workbook(fileName, new LoadOptions(LoadFormat.Excel97To2003));
                if (workbook != null)
                {
                    Worksheet worksheet = null;
                    if (sheetName != null)
                    {
                        worksheet = workbook.Worksheets[sheetName];
                    }
                    else
                    {
                        worksheet = workbook.Worksheets[0];
                    }
                    if (worksheet != null)
                    {
                        data = worksheet.Cells.ExportDataTable(0, 0, worksheet.Cells.MaxRow+1, worksheet.Cells.MaxColumn+1,
                            isFirstRowColumnName);

                        return data;
                    }
                }
                else
                {
                    return data;
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }

            return data;
        }
    }
}
```

## 流转byte[]
```csharp

public byte[] StreamToBytes(Stream stream) 

{

byte[] bytes = new byte[stream.Length]; 

stream.Read(bytes, 0, bytes.Length); 

// 设置当前流的位置为流的开始 

stream.Seek(0, SeekOrigin.Begin); 

return bytes; 

}


```

## 将byte[]转成 Stream
```csharp
/// 将 byte[] 转成 Stream

public Stream BytesToStream(byte[] bytes) 


{ 
Stream stream = new MemoryStream(bytes); 

return stream; 

}
```


## 防止死锁的运行异步方法
```csharp
var result = Task.Run(() => asyncGetValue()).Result;
Task.Run( () => asyncMethod()).Wait();
```

## 导出C#生成的公私钥
```csharp
 /// <summary>
        /// 导出公钥
        /// </summary>
        /// <param name="csp"></param>
        /// <param name="outputStream"></param>
        private static void ExportPublicKey(RSACryptoServiceProvider csp, TextWriter outputStream)
        {
            var parameters = csp.ExportParameters(false);
            using (var stream = new MemoryStream())
            {
                var writer = new BinaryWriter(stream);
                writer.Write((byte)0x30); // SEQUENCE
                using (var innerStream = new MemoryStream())
                {
                    var innerWriter = new BinaryWriter(innerStream);
                    innerWriter.Write((byte)0x30); // SEQUENCE
                    EncodeLength(innerWriter, 13);
                    innerWriter.Write((byte)0x06); // OBJECT IDENTIFIER
                    var rsaEncryptionOid = new byte[] { 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x01 };
                    EncodeLength(innerWriter, rsaEncryptionOid.Length);
                    innerWriter.Write(rsaEncryptionOid);
                    innerWriter.Write((byte)0x05); // NULL
                    EncodeLength(innerWriter, 0);
                    innerWriter.Write((byte)0x03); // BIT STRING
                    using (var bitStringStream = new MemoryStream())
                    {
                        var bitStringWriter = new BinaryWriter(bitStringStream);
                        bitStringWriter.Write((byte)0x00); // # of unused bits
                        bitStringWriter.Write((byte)0x30); // SEQUENCE
                        using (var paramsStream = new MemoryStream())
                        {
                            var paramsWriter = new BinaryWriter(paramsStream);
                            EncodeIntegerBigEndian(paramsWriter, parameters.Modulus); // Modulus
                            EncodeIntegerBigEndian(paramsWriter, parameters.Exponent); // Exponent
                            var paramsLength = (int)paramsStream.Length;
                            EncodeLength(bitStringWriter, paramsLength);
                            bitStringWriter.Write(paramsStream.GetBuffer(), 0, paramsLength);
                        }
                        var bitStringLength = (int)bitStringStream.Length;
                        EncodeLength(innerWriter, bitStringLength);
                        innerWriter.Write(bitStringStream.GetBuffer(), 0, bitStringLength);
                    }
                    var length = (int)innerStream.Length;
                    EncodeLength(writer, length);
                    writer.Write(innerStream.GetBuffer(), 0, length);
                }

                var base64 = Convert.ToBase64String(stream.GetBuffer(), 0, (int)stream.Length).ToCharArray();
                outputStream.WriteLine("-----BEGIN PUBLIC KEY-----");
                for (var i = 0; i < base64.Length; i += 64)
                {
                    outputStream.WriteLine(base64, i, Math.Min(64, base64.Length - i));
                }
                outputStream.WriteLine("-----END PUBLIC KEY-----");
            }
        }

        /// <summary>
        /// 导出私钥
        /// </summary>
        /// <param name="csp"></param>
        /// <param name="outputStream"></param>
        private static void ExportPrivateKey(RSACryptoServiceProvider csp, TextWriter outputStream)
        {
            if (csp.PublicOnly) throw new ArgumentException("CSP does not contain a private key", "csp");
            var parameters = csp.ExportParameters(true);
            using (var stream = new MemoryStream())
            {
                var writer = new BinaryWriter(stream);
                writer.Write((byte)0x30); // SEQUENCE
                using (var innerStream = new MemoryStream())
                {
                    var innerWriter = new BinaryWriter(innerStream);
                    EncodeIntegerBigEndian(innerWriter, new byte[] { 0x00 }); // Version
                    EncodeIntegerBigEndian(innerWriter, parameters.Modulus);
                    EncodeIntegerBigEndian(innerWriter, parameters.Exponent);
                    EncodeIntegerBigEndian(innerWriter, parameters.D);
                    EncodeIntegerBigEndian(innerWriter, parameters.P);
                    EncodeIntegerBigEndian(innerWriter, parameters.Q);
                    EncodeIntegerBigEndian(innerWriter, parameters.DP);
                    EncodeIntegerBigEndian(innerWriter, parameters.DQ);
                    EncodeIntegerBigEndian(innerWriter, parameters.InverseQ);
                    var length = (int)innerStream.Length;
                    EncodeLength(writer, length);
                    writer.Write(innerStream.GetBuffer(), 0, length);
                }

                var base64 = Convert.ToBase64String(stream.GetBuffer(), 0, (int)stream.Length).ToCharArray();
                outputStream.WriteLine("-----BEGIN RSA PRIVATE KEY-----");
                // Output as Base64 with lines chopped at 64 characters
                for (var i = 0; i < base64.Length; i += 64)
                {
                    
                    outputStream.WriteLine(base64, i, Math.Min(64, base64.Length - i));
                }
                outputStream.WriteLine("-----END RSA PRIVATE KEY-----");
                
            }
        }

//使用方法
/*
 RSACryptoServiceProvider rsa = new RSACryptoServiceProvider();
            rsa.FromXmlString(@"<RSAKeyValue><Modulus>qFZP5TVjcCmf7jEBHEnaum/yPq3WQ5aGPG35Z3Gt7a115cNkGM/wWZt1M4qQOlnYld37qx1QDIGXv61jChr30hBLVcG+YvBMKrA4CqI78M4t0spH0i4LM759fwdyWL22JgaobNZoBeMhBRAOtF0XT/LH/G+fd0sEjbX7nelAFL0=</Modulus><Exponent>AQAB</Exponent><P>xRHO1CAe2s7DizyI1wWXJPIPiocJIyYIe19DKZB6TWXgU85mhCFbwVx+iFOhn/NNzWCkl6xXOkUh3i/5OHX2/w==</P><Q>2qz2j4l1y9inWV/KWcYKqs8Uhjpvu8lwhQYwL2Oojo/5iCyMpQFzFV+dy2esjmhoCku/vX0JJy+S0waKMY+QQw==</Q><DP>MaCibVkJbCDVraK48y09OtiagVAwROG3ERqUV0tDAWq+a1x3BJ9B9BfO5ZXqBdXHqgjEak3ESbBPLxz1rfpHEQ==</DP><DQ>T82haX6j05GseQxhP2Przqwl9FptHl4ERzeb7B91ixl12kFPzoP56MntPycFrS7jESbVwaRY68kLzyFq221mGw==</DQ><InverseQ>WRnsQeILv0VX6fS1Q7XzHCo8a9gt9fZm5q24yxQZ5t5QWBMAJ/pKoDa3mEIiId9NYLYvek9cZN+brcWCNdprMA==</InverseQ><D>CC0A/mHkbXsoEFqC8kvH+swbGN46jNfPtzmkJlkIGIYXNsyRnP7kboW1YIZ3UM4yTb0VTw9CZwkYRK/4InKC3LXKrexT2F0NSl9qNbOlv8BtqOLL2qST9TLsWgzkhltriod8G2/F79GDSSZ0an8Uj3yHj5ChcaVwXhr0M8BAKMU=</D></RSAKeyValue>

");
            using (StreamWriter writer = File.CreateText("ExportPrivateKey.txt"))
            {
                ExportPrivateKey(rsa, writer);
            }
            rsa.FromXmlString(@"<RSAKeyValue><Modulus>qFZP5TVjcCmf7jEBHEnaum/yPq3WQ5aGPG35Z3Gt7a115cNkGM/wWZt1M4qQOlnYld37qx1QDIGXv61jChr30hBLVcG+YvBMKrA4CqI78M4t0spH0i4LM759fwdyWL22JgaobNZoBeMhBRAOtF0XT/LH/G+fd0sEjbX7nelAFL0=</Modulus><Exponent>AQAB</Exponent></RSAKeyValue>

");
            using (StreamWriter writer = File.CreateText("ExportPublicKey.txt"))
            {
                ExportPublicKey(rsa, writer);
            }

*/

```
参考链接:
1. [C# Export Private/Public RSA key from RSACryptoServiceProvider to PEM string](https://stackoverflow.com/questions/23734792/c-sharp-export-private-public-rsa-key-from-rsacryptoserviceprovider-to-pem-strin)
1. [通过 JSEncrypt 实现在线 JavaScript RSA 非对称加密解密](https://blog.zhengxianjun.com/online-tool/rsa/)



## 将时间字符串转为Json时间
```csharp
  /// <summary>
        /// 将时间字符串转为Json时间
        /// </summary>
        private static string ConvertDateStringToJsonDate(Match m)
        {
            string result = string.Empty;
            DateTime dt = DateTime.Parse(m.Groups[0].Value);
            dt = dt.ToUniversalTime();
            TimeSpan ts = dt - DateTime.Parse("1970-01-01");
            result = string.Format("\\/Date({0}+0800)\\/", ts.TotalMilliseconds);
            return result;
        }
```

## 
```csharp

```
