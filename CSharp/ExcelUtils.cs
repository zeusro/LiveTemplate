using System;
using System.Data;
using System.Data.OleDb;
using System.Data.SqlClient;
using System.IO;
using System.Collections;
using System.Configuration;
using System.Web;
using System.Collections.Generic;

namespace KingPoint.Commons
{
	/// <summary>
	/// ExcelUtils 的摘要说明。
	/// </summary>
	public static class ExcelUtils
	{
        /// <summary>
        /// 第一个表(sheet)的默认名称.
        /// </summary>
		public const string DEFAULT_FIRST_SHEET_NAME = "Sheet1";

        /// <summary>
        /// 读取 Excel 文件的连接字符串格式.
        /// </summary>
        public const string EXCEL_READ_CONN_STR = @"Provider=Microsoft.Jet.OLEDB.4.0;Data Source={0};Extended properties='Excel 8.0;HDR=Yes;IMEX=1;'";


        /// <summary>
        /// 写入 Excel 文件的连接字符串格式.
        /// </summary>
		public const string EXCEL_WRITE_CONN_STR = @"Provider=Microsoft.Jet.OLEDB.4.0;Data Source={0};Extended properties='Excel 8.0;'";

        /// <summary>
        /// 选择 Excel 数据的表达式字符串.
        /// </summary>
		public const string EXCEL_SELECT_STR = "select * from [{0}$]";

        /// <summary>
        /// 选择 Excel 数据的表达式字符串.
        /// </summary>
        public const string EXCEL_SELECT_STR2 = "select * from [{0}$A:IU]";

		/// <summary>
		/// 读取excel文件.
		/// </summary>
		/// <param name="fileName">文件名.</param>
		/// <param name="sheetName">工作表名.</param>
		/// <returns></returns>
		public static DataSet ReadExcel(string fileName, string sheetName)
		{
            return ReadExcel(fileName, sheetName, false);
		}

        /// <summary>
        /// 读取excel文件.
        /// </summary>
        /// <param name="fileName">文件名.</param>
        /// <param name="sheetName">工作表名.</param>
        /// <param name="deleteFileAfterRead">是否在读取之后删除文件.</param>
        /// <returns></returns>
        public static DataSet ReadExcel(string fileName, string sheetName, bool deleteFileAfterRead)
        {
            string strConn = String.Format(EXCEL_READ_CONN_STR, fileName);
            string strSelect = String.Format(EXCEL_SELECT_STR, sheetName);
            OleDbConnection oleDbConnection = new OleDbConnection(strConn);
            OleDbDataAdapter oleDbDataAdapter = new OleDbDataAdapter(strSelect, oleDbConnection);
            DataSet ds = new DataSet();
            oleDbDataAdapter.Fill(ds, sheetName);
            if (deleteFileAfterRead)
            {
                File.Delete(fileName);
            }            
            return ds;
        }
        public static DataSet ReadExcel2(string fileName, string sheetName, bool deleteFileAfterRead)
        {
            string strConn = String.Format(EXCEL_READ_CONN_STR, fileName);
            string strSelect = String.Format(EXCEL_SELECT_STR2, sheetName);
            OleDbConnection oleDbConnection = new OleDbConnection(strConn);
            OleDbDataAdapter oleDbDataAdapter = new OleDbDataAdapter(strSelect, oleDbConnection);
            DataSet ds = new DataSet();
            oleDbDataAdapter.Fill(ds, sheetName);
            if (deleteFileAfterRead)
            {
                File.Delete(fileName);
            }
            return ds;
        }
		/// <summary>
		/// 读取excel文件.
		/// </summary>
		/// <param name="fileName">文件名.</param>
		/// <param name="sheetNames">多个工作表.</param>
		/// <returns></returns>
		public static DataSet ReadExcel(string fileName, string[] sheetNames)
		{
            return ReadExcel(fileName, sheetNames, false);
		}

        /// <summary>
        /// 读取excel文件.
        /// </summary>
        /// <param name="fileName">文件名.</param>
        /// <param name="sheetNames">多个工作表.</param>
        /// <param name="deleteFileAfterRead">是否在读取之后删除文件.</param>
        /// <returns></returns>
        public static DataSet ReadExcel(string fileName, string[] sheetNames, bool deleteFileAfterRead)
        {
            OleDbConnection oleDbConnection = new OleDbConnection(String.Format(EXCEL_READ_CONN_STR, fileName));
            OleDbDataAdapter oleDbDataAdapter = new OleDbDataAdapter("", oleDbConnection);
            DataSet ds = new DataSet();
            for (int i = 0; i < sheetNames.Length; i++)
            {
                oleDbDataAdapter.SelectCommand.CommandText = String.Format(EXCEL_SELECT_STR, sheetNames[i]);
                oleDbDataAdapter.Fill(ds, sheetNames[i]);
            }
            if (deleteFileAfterRead)
            {
                File.Delete(fileName);
            }   
            return ds;
        }

		/// <summary>
		/// 获取excel的sheet名.
		/// </summary>
		/// <param name="fileName"></param>
		/// <returns></returns>
		public static DataTable ReadExcelSheetNames(string fileName)
		{
			OleDbConnection oleDbConnection = new OleDbConnection(String.Format(EXCEL_READ_CONN_STR, fileName));
			oleDbConnection.Open();
			DataTable dt = oleDbConnection.GetOleDbSchemaTable(OleDbSchemaGuid.Tables,null);
			oleDbConnection.Close();
			return dt;
		}

		/// <summary>
		/// 导入EXCEL至Sql数据库.
		/// </summary>
		/// <param name="srcDt"></param>
		/// <param name="strConn"></param>
		/// <param name="destTableName"></param>
		/// <param name="refFieldNames"></param>
		/// <returns></returns>
		public static int ImportExcel2DB(DataTable srcDt,string strConn,string destTableName,string[,] refFieldNames)
		{
			string strSelect= string.Format("select top 1 * from [{0}]", destTableName);
			SqlConnection conn = new SqlConnection(strConn);
			SqlDataAdapter da = new SqlDataAdapter(strSelect, conn);
			DataSet ds = new DataSet();
			da.Fill(ds, destTableName);
			DataTable destDt=ds.Tables[0];

			//获取源数据中的字段名称与数据库中的字段名称
			int fieldLength = refFieldNames.GetLength(0);
			string[] srcFieldNames=new string[fieldLength];
			string[] destFieldNames=new string[fieldLength];

			//转换标致：0-不转换,1-数值型,2-日期型
			string[] convertFlags=new string[fieldLength];
			
			for (int i=0;i<fieldLength;i++)
			{
				srcFieldNames[i] = refFieldNames.GetValue(i, 0).ToString();
				destFieldNames[i] = refFieldNames.GetValue(i, 1).ToString();
				convertFlags[i] = refFieldNames.GetValue(i, 2).ToString();
			}

			foreach(DataRow drSrc in srcDt.Rows)
			{
				DataRow  drNew = destDt.NewRow();

				for (int i=0;i<fieldLength;i++)
				{
					if (convertFlags[i]=="1")
						drNew[destFieldNames[i]] = int.Parse(drSrc[srcFieldNames[i]].ToString());
					else
					{
						if (convertFlags[i]=="2")
							drNew[destFieldNames[i]] = DateTime.Parse(drSrc[srcFieldNames[i]].ToString());
						else
							drNew[destFieldNames[i]] = drSrc[srcFieldNames[i]];
					}
				}

				destDt.Rows.Add(drNew);
			}

			SqlCommandBuilder cb  = new SqlCommandBuilder(da);
			cb.QuotePrefix = "[";
			cb.QuoteSuffix = "]";

			return da.Update(ds,destTableName);
		}

		/// <summary>
		/// 导出到Excel.
		/// </summary>
		/// <param name="dsSrc"></param>
		/// <param name="tableName"></param>
		/// <param name="templateExcelFileName"></param>
		/// <param name="sheetName"></param>
		/// <param name="destExcelFileName"></param>
		public static void Export2Excel(DataSet dsSrc, string tableName, string templateExcelFileName, string sheetName, string destExcelFileName)
		{
            Export2Excel(dsSrc.Tables[tableName], templateExcelFileName, sheetName, destExcelFileName);
		}

        /// <summary>
        /// 导出到Excel.
        /// </summary>
        /// <param name="dtSrc"></param>
        /// <param name="templateExcelFileName"></param>
        /// <param name="sheetName"></param>
        /// <param name="destExcelFileName"></param>
        public static void Export2Excel(DataTable dtSrc, string templateExcelFileName, string sheetName, string destExcelFileName)
        {
            if (!System.IO.File.Exists(templateExcelFileName))
            {
                throw new ArgumentException(String.Format("文件{0}不存在", templateExcelFileName), "templateExcelFileName");
            }

            System.IO.File.Copy(templateExcelFileName, destExcelFileName, true);

            string strConn = String.Format(EXCEL_WRITE_CONN_STR, destExcelFileName);
            string strSelect = String.Format(EXCEL_SELECT_STR, sheetName);
            OleDbConnection conn = new OleDbConnection(strConn);
            OleDbDataAdapter da = new OleDbDataAdapter(strSelect, conn);
            DataSet ds = new DataSet();
            da.Fill(ds, sheetName);

            foreach (DataRow drSrc in dtSrc.Rows)
            {
                DataRow drNew = ds.Tables[sheetName].NewRow();
                foreach (DataColumn dc in dtSrc.Columns)
                {
                    drNew[dc.ColumnName] = drSrc[dc.ColumnName];
                }
                ds.Tables[sheetName].Rows.Add(drNew);
            }

            OleDbCommandBuilder cb = new OleDbCommandBuilder(da);
            cb.QuotePrefix = "[";
            cb.QuoteSuffix = "]";
            da.Update(ds, sheetName);
        }

        /// <summary>
        /// 导出数据到Excel(有对应的Excel模板文件).
        /// </summary>
        /// <param name="dsSrc">源数据集.</param>
        /// <param name="tableName">要导出的表.</param>
        /// <param name="templateExcelFileName">对应的Excel模板文件.</param>
        /// <param name="sheetName">工作簿名.</param>
        /// <param name="refFieldNames">tableName表中字段与sheetName工作簿字段的对应二维数组.</param>
        /// <param name="destExcelFileName">目标Excel.</param>
		public static void Export2Excel(DataSet dsSrc, string tableName, string templateExcelFileName, string sheetName, string[,] refFieldNames, string destExcelFileName)
		{
			int fieldLength = refFieldNames.GetLength(0);
			string[] srcFieldNames = new string[fieldLength];
			string[] destFieldNames = new string[fieldLength];
			for (int i=0; i<fieldLength; i++)
			{
				srcFieldNames[i] = refFieldNames.GetValue(i, 0).ToString();
				destFieldNames[i] = refFieldNames.GetValue(i, 1).ToString();
			}

            Export2Excel(dsSrc.Tables[tableName], templateExcelFileName, sheetName, srcFieldNames, destFieldNames, destExcelFileName);
		}

		/// <summary>
		/// 导出数据到Excel(有对应的Excel模板文件).
		/// </summary>
		/// <param name="dsSrc">源数据集</param>
		/// <param name="tableName">要导出的表</param>
		/// <param name="templateExcelFileName">对应的Excel模板文件</param>
		/// <param name="sheetName">工作簿名</param>
		/// <param name="srcFieldNames">数据表中字段</param>
		/// <param name="destFieldNames">对应Excel的字段</param>
		/// <param name="destExcelFileName">目标Excel</param>
		public static void Export2Excel(DataSet dsSrc, string tableName, string templateExcelFileName, string sheetName, string[] srcFieldNames, string[] destFieldNames, string destExcelFileName)
		{
            Export2Excel(dsSrc.Tables[tableName], templateExcelFileName, sheetName, srcFieldNames, destFieldNames, destExcelFileName);
		}

        /// <summary>
        /// 导出数据到Excel(有对应的Excel模板文件)
        /// </summary>
        /// <param name="dtSrc">源数据表</param>
        /// <param name="templateExcelFileName">对应的Excel模板文件</param>
        /// <param name="sheetName">工作簿名</param>
        /// <param name="srcFieldNames">数据表中字段</param>
        /// <param name="destFieldNames">对应Excel的字</param>
        /// <param name="destExcelFileName">目标Excel</param>
        public static void Export2Excel(DataTable dtSrc, string templateExcelFileName, string sheetName, string[] srcFieldNames, string[] destFieldNames, string destExcelFileName)
        {
            if (null == srcFieldNames
                || null == destFieldNames
                || srcFieldNames.Length != destFieldNames.Length
                || srcFieldNames.Length <= 0)
            {
                throw new ArgumentException("srcFieldNames/destFieldNames (null/empty/length) is invalid");
            }

            if (!System.IO.File.Exists(templateExcelFileName))
            {
                throw new ArgumentException(String.Format("文件{0}不存在", templateExcelFileName), "templateExcelFileName");
            }

            System.IO.File.Copy(templateExcelFileName, destExcelFileName, true);

            string strConn = String.Format(EXCEL_WRITE_CONN_STR, destExcelFileName);
            string strSelect = String.Format(EXCEL_SELECT_STR, sheetName);
            OleDbConnection conn = new OleDbConnection(strConn);
            OleDbDataAdapter da = new OleDbDataAdapter(strSelect, conn);
            DataSet ds = new DataSet();
            da.Fill(ds, sheetName);

            int fieldLength = srcFieldNames.Length;

            foreach (DataRow drSrc in dtSrc.Rows)
            {
                DataRow drNew = ds.Tables[sheetName].NewRow();
                for (int i = 0; i < fieldLength; i++)
                {
                    drNew[destFieldNames[i]] = drSrc[srcFieldNames[i]];
                }
                ds.Tables[sheetName].Rows.Add(drNew);
            }

            OleDbCommandBuilder cb = new OleDbCommandBuilder(da);
            cb.QuotePrefix = "[";
            cb.QuoteSuffix = "]";
            da.Update(ds, sheetName);
        }
		

        /// <summary>
        /// 将列表数据导入到Excel中
        /// </summary>
        /// <param name="templateExcelFileName"></param>
        /// <param name="fieldNames"></param>
        /// <param name="datas"></param>
        /// <param name="destExcelFileName"></param>
        public static void Export2Excel(string templateExcelFileName, string[] fieldNames, IList<string[]> datas, string destExcelFileName)
        {
            Export2Excel(templateExcelFileName, DEFAULT_FIRST_SHEET_NAME, fieldNames, datas, destExcelFileName);
        }

        /// <summary>
        /// 将列表数据导入到Excel中
        /// </summary>
        /// <param name="templateExcelFileName"></param>
        /// <param name="sheetName"></param>
        /// <param name="fieldNames"></param>
        /// <param name="datas"></param>
        /// <param name="destExcelFileName"></param>
        public static void Export2Excel(string templateExcelFileName, string sheetName, string[] fieldNames, IList<string[]> datas, string destExcelFileName)
        {
            if (null == fieldNames
                || fieldNames.Length == 0
                || null == datas
                || datas.Count == 0)
            {
                return;
            }

            if (fieldNames.Length > datas[0].Length)
            {
                throw new ArgumentOutOfRangeException("fieldNames", "数据的字段数不能小于模板的字段数");
            }

            System.IO.File.Copy(templateExcelFileName, destExcelFileName, true);

            string strConn = String.Format(EXCEL_WRITE_CONN_STR, destExcelFileName);
            OleDbConnection conn = new OleDbConnection(strConn);
            OleDbCommand cmd = new OleDbCommand();
            cmd.Connection = conn;

            int fieldLength = fieldNames.Length;
            string[] arrParms = new string[fieldLength];
            for (int i = 0; i < fieldLength; i++)
            {
                arrParms[i] = "@" + i.ToString();

                OleDbParameter pa = new OleDbParameter();
                pa.ParameterName = arrParms[i];
                cmd.Parameters.Add(pa);
            }
            string fields = StringUtils.Join(fieldNames, "[", "]", ",");
            string parms = StringUtils.Join(arrParms, ",");
            cmd.CommandText = String.Format("insert into [{0}$] ({1}) values ({2})", sheetName, fields, parms);

            using (conn)
            {
                conn.Open();
                using (OleDbTransaction trans = conn.BeginTransaction())
                {
                    cmd.Transaction = trans;
                    try
                    {
                        for (int i = 0, icount = datas.Count; i < icount; i++)
                        {
                            string[] arr = datas[i];
                            for (int j = 0; j < fieldLength; j++)
                            {
                                cmd.Parameters[arrParms[j]].Value = arr[j];
                            }
                            cmd.ExecuteNonQuery();
                        }
                        trans.Commit();
                    }
                    catch (SqlException ex)
                    {
                        if (null != trans)
                        {
                            trans.Rollback();
                        }
                        throw ex;
                    }
                }
            }            
        }

        /// <summary>
        /// 将WebControl数据列表导出至Excel
        /// </summary>
        /// <param name="objControl"></param>
        /// <param name="strFileName"></param>
        public static void Export2Excel(System.Web.UI.Control objControl, string strFileName)
        {
            strFileName = System.Web.HttpUtility.UrlEncode(strFileName, System.Text.Encoding.UTF8);

            System.Web.HttpContext.Current.Response.Clear();
            System.Web.HttpContext.Current.Response.Buffer = true;
            System.Web.HttpContext.Current.Response.Charset = "GB2312";
            System.Web.HttpContext.Current.Response.AppendHeader("Content-Disposition", "attachment; filename=" + strFileName + ".xls");
            System.Web.HttpContext.Current.Response.ContentEncoding = System.Text.Encoding.UTF7;
            System.Web.HttpContext.Current.Response.ContentType = "application/ms-excel";

            System.Globalization.CultureInfo myCItrad = new System.Globalization.CultureInfo("ZH-CN", true);
            System.IO.StringWriter oStringWriter = new System.IO.StringWriter(myCItrad);
            System.Web.UI.HtmlTextWriter oHtmlTextWriter = new System.Web.UI.HtmlTextWriter(oStringWriter);

            objControl.RenderControl(oHtmlTextWriter);

            objControl = null;
            System.Web.HttpContext.Current.Response.Write(oStringWriter.ToString());
            System.Web.HttpContext.Current.Response.Buffer = false;
            System.Web.HttpContext.Current.Response.End();
        }

        /// <summary>
        /// 导出数据到Excel(有对应的Excel模板文件)
        /// </summary>
        /// <param name="dtSrc">源数据表</param>
        /// <param name="templateExcelFileName">对应的Excel模板文件</param>
        /// <param name="sheetName">工作簿名</param>
        /// <param name="srcFieldNames">数据表中字段</param>
        /// <param name="destExcelFileName">目标Excel</param>
        public static void Export2Excel(DataTable dtSrc, string templateExcelFileName, string sheetName, string[] srcFieldNames, string destExcelFileName)
        {
            if (null == srcFieldNames
                || srcFieldNames.Length <= 0)
            {
                throw new ArgumentException("srcFieldNames/destFieldNames (null/empty/length) is invalid");
            }

            if (!System.IO.File.Exists(templateExcelFileName))
            {
                throw new ArgumentException("templateExcelFileName", templateExcelFileName);
            }

            System.IO.File.Copy(templateExcelFileName, destExcelFileName, true);

            string strConn = String.Format(EXCEL_WRITE_CONN_STR, destExcelFileName);
            string strSelect = String.Format(EXCEL_SELECT_STR, sheetName);
            OleDbConnection conn = new OleDbConnection(strConn);
            OleDbDataAdapter da = new OleDbDataAdapter(strSelect, conn);
            DataSet ds = new DataSet();
            da.Fill(ds, sheetName);

            int fieldLength = srcFieldNames.Length;

            foreach (DataRow drSrc in dtSrc.Rows)
            {
                DataRow drNew = ds.Tables[sheetName].NewRow();
                for (int i = 0; i < fieldLength; i++)
                {
                    drNew[i] = drSrc[srcFieldNames[i]];
                }
                ds.Tables[sheetName].Rows.Add(drNew);
            }

            OleDbCommandBuilder cb = new OleDbCommandBuilder(da);
            cb.QuotePrefix = "[";
            cb.QuoteSuffix = "]";
            da.Update(ds, sheetName);
        } 


		/// <summary>
		/// 判断是否为excel文件
		/// </summary>
		/// <param name="fileName"></param>
		/// <returns></returns>
		public static bool IsExcelFile(string fileName)
		{
			int idx = fileName.LastIndexOf(".");
			string fixName = fileName.Substring(idx + 1);
			return (fixName.ToLower() == "xls");
		}
	}
}
