using System;

namespace Common.Tools
{
    /// <summary>
    /// 时间工具
    /// </summary>
    public static class TimeUtil
    {
        //unix 时间戳开始时间
        private static readonly DateTime UnixDateTime = new DateTime(1970, 1, 1, 0, 0, 0, 0);

        //服务器unix 时间戳
        public static long ServerUnixMillisecond { get; set; }

        /// <summary>
        /// 获取时间戳 毫秒
        /// </summary>
        /// <returns></returns>
        public static long CurrentTimeMillis()
        {
            TimeSpan ts = DateTime.Now - UnixDateTime;
            return Convert.ToInt64(ts.TotalMilliseconds);
        }

        /// <summary>
        /// 当前格式化时间  yyyy-MM-dd HH:mm:ss
        /// </summary>
        /// <returns></returns>
        public static string CurrentFormatTime()
        {
            return DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss");
        }
    }
}