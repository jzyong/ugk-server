//#define NeedTry

using System;
using System.Collections.Generic;
using UnityEngine;

namespace Common.Tools
{
    /// <summary>
    /// 事件控制器
    /// </summary>
    /// <typeparam name="Type"></typeparam>
    public class BaseEventManager<Type>
    {
        Dictionary<Type, Delegate> eventMap;

        public BaseEventManager()
        {
            eventMap = new Dictionary<Type, Delegate>();
        }

        /// <summary>
        /// 添加事件
        /// </summary>
        /// <param name="eventName"></param>
        /// <param name="eventHandle"></param>
        public void AddEvent(Type eventName, System.Action eventHandle)
        {
            if (!eventMap.ContainsKey(eventName))
                eventMap[eventName] = (System.Action) eventHandle;
            else
            {
                try
                {
                    eventMap[eventName] = Delegate.Combine((System.Action) eventMap[eventName], eventHandle);
                }
                catch
                {
                    Debug.LogError($"添加事件：{eventName} 失败");
                }
            }
        }

        public void AddEvent<T>(Type eventName, Action<T> eventHandle)
        {
            if (!eventMap.ContainsKey(eventName))
                eventMap[eventName] = eventHandle;
            else
                eventMap[eventName] = (Action<T>) Delegate.Combine((Action<T>) eventMap[eventName], eventHandle);
        }

        public void AddEvent<T, U>(Type eventName, Action<T, U> eventHandle)
        {
            if (!eventMap.ContainsKey(eventName))
                eventMap[eventName] = eventHandle;
            else
                eventMap[eventName] = Delegate.Combine(eventMap[eventName], eventHandle);
        }

        public void AddEvent<T, U, V>(Type eventName, Action<T, U, V> eventHandle)
        {
            if (!eventMap.ContainsKey(eventName))
                eventMap[eventName] = eventHandle;
            else
                eventMap[eventName] = Delegate.Combine(eventMap[eventName], eventHandle);
        }

        public void AddEvent<T, U, V, W>(Type eventName, Action<T, U, V, W> eventHandle)
        {
            if (!eventMap.ContainsKey(eventName))
                eventMap[eventName] = eventHandle;
            else
                eventMap[eventName] = Delegate.Combine(eventMap[eventName], eventHandle);
        }

        public void OnEvent(Type eventName)
        {
            if (eventMap.ContainsKey(eventName))
            {
                Delegate gata = eventMap[eventName];
                if (gata == null)
                    return;
                Delegate[] gateList = gata.GetInvocationList();
                int count = gateList.Length;
                for (int i = 0; i < count; i++)
                {
                    string str = gateList[i].Target.ToString();
                    if (str == "null")
                    {
                        eventMap[eventName] = gata = Delegate.Remove(gata, gateList[i]);
                        continue;
                    }

                    System.Action action = gateList[i] as System.Action;
                    if (action != null)
                    {
                        action();
                    }
                    else
                    {
                        Debug.Log(string.Format("对象{0}脚本注册的GameEventEnum.{1}事件参数类型未转换成功", str, eventName));
                    }
                }
            }
        }

        public void OnEvent<T>(Type eventName, T arg1)
        {
            if (eventMap.ContainsKey(eventName))
            {
                Delegate gata = eventMap[eventName];
                if (gata == null)
                    return;
                Delegate[] gateList = gata.GetInvocationList();
                int count = gateList.Length;
                for (int i = 0; i < count; i++)
                {
                    string str = gateList[i].Target.ToString();
                    if (str == "null")
                    {
                        eventMap[eventName] = gata = Delegate.Remove(gata, gateList[i]);

                        continue;
                    }

                    Action<T> action = gateList[i] as Action<T>;
                    if (action != null)
                    {
                        action(arg1);
                    }
                    else
                    {
                        Debug.Log($"对象{str}脚本注册的GameEventEnum.{eventName}事件参数类型未转换成功");
                    }
                }
            }
        }

        public void OnEvent<T, U>(Type eventName, T arg1, U arg2)
        {
            if (eventMap.ContainsKey(eventName))
            {
                Delegate gata = eventMap[eventName];
                if (gata == null)
                    return;
                Delegate[] gateList = gata.GetInvocationList();
                int count = gateList.Length;
                for (int i = 0; i < count; i++)
                {
                    string str = gateList[i].Target.ToString();
                    if (str == "null")
                    {
                        eventMap[eventName] = gata = Delegate.Remove(gata, gateList[i]);
                        continue;
                    }

                    Action<T, U> action = gateList[i] as Action<T, U>;
                    if (action != null)
                    {
                        action(arg1, arg2);
                    }
                    else
                    {
                        UnityEngine.Debug.Log(string.Format("对象{0}脚本注册的GameEventEnum.{1}事件参数类型未转换成功", str, eventName));
                    }
                }
            }
        }

        public void OnEvent<T, U, V>(Type eventName, T arg1, U arg2, V arg3)
        {
            if (eventMap.ContainsKey(eventName))
            {
                Delegate gata = eventMap[eventName];
                if (gata == null)
                    return;
                Delegate[] gateList = gata.GetInvocationList();
                int count = gateList.Length;
                for (int i = 0; i < count; i++)
                {
                    string str = gateList[i].Target.ToString();
                    if (str == "null")
                    {
                        eventMap[eventName] = gata = Delegate.Remove(gata, gateList[i]);
                        continue;
                    }

                    Action<T, U, V> action = gateList[i] as Action<T, U, V>;
                    if (action != null)
                    {
                        action(arg1, arg2, arg3);
                    }
                    else
                    {
                        UnityEngine.Debug.Log(string.Format("对象{0}脚本注册的GameEventEnum.{1}事件参数类型未转换成功", str, eventName));
                    }
                }
            }
        }

        public void OnEvent<T, U, V, W>(Type eventName, T arg1, U arg2, V arg3, W arg4)
        {
            if (eventMap.ContainsKey(eventName))
            {
                Delegate gata = eventMap[eventName];
                if (gata == null)
                    return;
                Delegate[] gateList = gata.GetInvocationList();
                int count = gateList.Length;
                for (int i = 0; i < count; i++)
                {
                    string str = gateList[i].Target.ToString();
                    if (str == "null")
                    {
                        eventMap[eventName] = gata = Delegate.Remove(gata, gateList[i]);
                        continue;
                    }

                    Action<T, U, V, W> action = gateList[i] as Action<T, U, V, W>;
                    if (action != null)
                    {
                        action(arg1, arg2, arg3, arg4);
                    }
                    else
                    {
                        UnityEngine.Debug.Log(string.Format("对象{0}脚本注册的GameEventEnum.{1}事件参数类型未转换成功", str, eventName));
                    }
                }
            }
        }

        public void RemoveEvent(Type eventName)
        {
            if (eventMap.ContainsKey(eventName))
            {
                eventMap.Remove(eventName);
            }
        }

        public void RemoveEvent(Type eventName, System.Action eventHandle)
        {
            Delegate gate = null;
            if (eventMap.TryGetValue(eventName, out gate))
                eventMap[eventName] = Delegate.Remove(gate, eventHandle);
        }

        public void RemoveEvent<T>(Type eventName, Action<T> eventHandle)
        {
            Delegate gate = null;
            if (eventMap.TryGetValue(eventName, out gate))
                eventMap[eventName] = Delegate.Remove(gate, eventHandle);
        }

        public void RemoveEvent<T, U>(Type eventName, Action<T, U> eventHandle)
        {
            Delegate gate = null;
            if (eventMap.TryGetValue(eventName, out gate))
                eventMap[eventName] = Delegate.Remove(gate, eventHandle);
        }

        public void RemoveEvent<T, U, V>(Type eventName, Action<T, U, V> eventHandle)
        {
            Delegate gate = null;
            if (eventMap.TryGetValue(eventName, out gate))
                eventMap[eventName] = Delegate.Remove(gate, eventHandle);
        }

        public void RemoveEvent<T, U, V, W>(Type eventName, Action<T, U, V, W> eventHandle)
        {
            Delegate gate = null;
            if (eventMap.TryGetValue(eventName, out gate))
                eventMap[eventName] = Delegate.Remove(gate, eventHandle);
        }

        public void Clear()
        {
            eventMap.Clear();
        }
    }
}