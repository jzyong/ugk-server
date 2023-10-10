using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Threading;

namespace Common.Tools
{
    /**
     * 自定义日志写文件
     */
    public class Log
    {
        //日志存储路径
        public static string[] LogPaths = { "./log" };

        //写数据级别
        public static LogLevel WriteLevel = LogLevel.Close;

        //写日志线程
        private static Thread writeThread;

        //日志缓存队列
        private static Queue<string> logQueue;

        /// <summary>
        /// 日志级别
        /// </summary>
        public enum LogLevel
        {
            Debug = 1,
            Info = 2,
            Warn = 3,
            Error = 4,
            Close = 100,
        }

        /// <summary>
        /// 格式化消息
        /// </summary>
        /// <param name="msg"></param>
        /// <returns></returns>
        private static string FormatMessage(string msg, LogLevel logLevel)
        {
            // 注意：StackTrace的参数1非常重要，表示获得父一级函数调用的相关信息。
            // 如果修改为0，则返回本身的行列信息，即TraceMethodInfo()函数的信息。
            StackFrame st = new StackTrace(2, true).GetFrame(0);
            var fileName=st.GetFileName();
            if (fileName!=null)
            {
                var strings = fileName.Split("\\");
                fileName = strings[strings.Length - 1];
            }
            
            return
                $"{TimeUtil.CurrentFormatTime()} [{logLevel.ToString()}]{fileName}({st.GetFileLineNumber()})--> {msg}";
        }

        public static void Debug(string msg)
        {
            msg = FormatMessage(msg, LogLevel.Debug);
            UnityEngine.Debug.Log(msg);
            if (WriteLevel <= LogLevel.Debug)
            {
                WriteLog(msg);
            }
        }

        public static void Info(string msg)
        {
            msg = FormatMessage(msg, LogLevel.Info);
            UnityEngine.Debug.Log(msg);
            if (WriteLevel <= LogLevel.Info)
            {
                WriteLog(msg);
            }
        }

        public static void Warn(string msg)
        {
            msg = FormatMessage(msg, LogLevel.Warn);
            UnityEngine.Debug.Log(msg);
            if (WriteLevel <= LogLevel.Warn)
            {
                WriteLog(msg);
            }
        }

        public static void Error(string msg)
        {
            msg = FormatMessage(msg, LogLevel.Error);
            UnityEngine.Debug.Log(msg);
            if (WriteLevel <= LogLevel.Error)
            {
                WriteLog(msg);
            }
        }


        /**
         * 输出通用日志文件
         */
        public static void WriteLog(string msg)
        {
            if (WriteLevel == LogLevel.Close)
            {
                return;
            }

            if (logQueue != null)
            {
                logQueue.Enqueue(msg);
            }
            else
            {
                logQueue = new Queue<string>(100);
                logQueue.Enqueue(msg);

                string sFilePath = LogPaths[0];
                string sFileName = "game_"+DateTime.Now.ToString("yyyy-MM-dd") + ".log";
                //文件的绝对路径
                sFileName = Path.Combine(sFilePath, sFileName);
                //验证路径是否存在,不存在则创建
                if (!Directory.Exists(LogPaths[0]))
                {
                    Directory.CreateDirectory(sFilePath);
                }


                writeThread = new Thread((() =>
                {
                    while (true)
                    {
                        if (logQueue.Count < 1)
                        {
                            Thread.Sleep(100);
                            continue;
                        }

                        FileStream fs;
                        StreamWriter sw;
                        //验证文件是否存在，有则追加，无则创建
                        fs = File.Exists(sFileName)
                            ? new FileStream(sFileName, FileMode.Append, FileAccess.Write)
                            : new FileStream(sFileName, FileMode.Create, FileAccess.Write);
                        sw = new StreamWriter(fs);
                        //日志内容
                        sw.WriteLine(logQueue.Dequeue());
                        sw.Close();
                        fs.Close();
                    }
                }));
                writeThread.Start();
            }
        }
    }
}