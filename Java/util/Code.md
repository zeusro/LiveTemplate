
1. []()
1. []()
1. []()
1. [格式化为数据库支持的时间格式](#格式化为数据库支持的时间格式)
1. [对json字符串格式化输出](#对json字符串格式化输出)
1. [列举可用数据驱动](#列举可用数据驱动)
1. [复制文本到剪贴板](#复制文本到剪贴板)

 
## 

```java

```


## 格式化为数据库支持的时间格式

```java
final java.text.SimpleDateFormat sdf = new java.text.SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
sdf.format(new Date());
```

## 对json字符串格式化输出
 ```java
    /**
     * 对json字符串格式化输出
     *
     * @param jsonStr JSON文本
     * @return 格式化过的JSON
     */
    public static String formatJson(String jsonStr) {
        if (null == jsonStr || "".equals(jsonStr)) {
            return "";
        }
        StringBuilder sb = new StringBuilder();
        char last = '\0';
        char current = '\0';
        int indent = 0;
        for (int i = 0; i < jsonStr.length(); i++) {
            last = current;
            current = jsonStr.charAt(i);
            switch (current) {
                case '{':
                case '[':
                    sb.append(current);
                    sb.append('\n');
                    indent++;
                    addIndentBlank(sb, indent);
                    break;
                case '}':
                case ']':
                    sb.append('\n');
                    indent--;
                    addIndentBlank(sb, indent);
                    sb.append(current);
                    break;
                case ',':
                    sb.append(current);
                    if (last != '\\') {
                        sb.append('\n');
                        addIndentBlank(sb, indent);
                    }
                    break;
                default:
                    sb.append(current);
            }
        }

        return sb.toString();
    }

    /**
     * 添加space
     *
     * @param sb     sb
     * @param indent 缩进
     */
    private static void addIndentBlank(StringBuilder sb, int indent) {
        for (int i = 0; i < indent; i++) {
            sb.append('\t');
        }
    }
```

## 列举可用数据驱动

```java
import java.sql.Driver;
import java.sql.DriverManager;
import java.util.Enumeration;

    private static void listDrivers() {
        Enumeration driverList = DriverManager.getDrivers();
        while (driverList.hasMoreElements()) {
            Driver driverClass = (Driver) driverList.nextElement();
            System.out.println(" dirver:" + driverClass.getClass().getName());
        }
    }
```

## 复制文本到剪贴板

```java
private void copyToClipboard(String content) {
        Clipboard clipboard = Toolkit.getDefaultToolkit().getSystemClipboard();
        Transferable trandata = new StringSelection(content);
        clipboard.setContents(trandata, null);
        out.println("responseJson已复制到剪贴板");
    }
```
