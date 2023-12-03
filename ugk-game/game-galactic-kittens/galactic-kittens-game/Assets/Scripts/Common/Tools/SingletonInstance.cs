using UnityEngine;

namespace Common.Tools
{
    /// <summary>
    /// 持久单例,不需要继承MonoBehaviour的，必须有无参构造方法的引用类型
    /// </summary>
    public class SingletonInstance<T> where T : class, new()
    {
        private static T instance = default(T);


        public static T Instance
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
    
    /// <summary>
    /// 继承 MonoBehaviour的单例
    /// </summary>
    /// <typeparam name="T"></typeparam>
    public class Singleton<T> : MonoBehaviour where T : Component
    {
        public static T Instance { get; private set; }

        public virtual void Awake()
        {
            if (Instance == null)
            {
                Instance = this as T;
            }
            else
            {
                Destroy(gameObject);
            }
        }
    }

    /// <summary>
    /// 继承 MonoBehaviour的单例且不销毁
    /// </summary>
    /// <typeparam name="T"></typeparam>
    public class SingletonPersistent<T> : MonoBehaviour where T : Component
    {
        public static T Instance { get; private set; }

        public virtual void Awake()
        {
            if (Instance == null)
            {
                Instance = this as T;
                DontDestroyOnLoad(this);
            }
            else
            {
                Destroy(gameObject);
            }
        }
    }
}