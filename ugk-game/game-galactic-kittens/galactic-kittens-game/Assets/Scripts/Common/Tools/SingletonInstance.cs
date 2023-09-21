namespace Common.Tools
{
    /// <summary>
    /// 持久单例,不需要继承MonoBehaviour的，必须有无参构造方法的引用类型
    /// </summary>
    public class SingletonInstance<T> where T : class, new()
    {
        private static T instance = default(T);


        public static T Singleton
        {
            get
            {
                if (instance == null)
                {
                    instance = new T();
                }

                return instance;
            }
            set => instance = value;
        }
    }
}