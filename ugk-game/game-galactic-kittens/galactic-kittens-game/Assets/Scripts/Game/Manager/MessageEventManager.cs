using Common.Tools;

namespace Game.Manager
{


    /// <summary>
    /// 消息事件枚举
    /// </summary>
    public enum MessageEvent
    {
        Login,
    }
    
    /// <summary>
    /// 消息事件处理
    /// </summary>
    public class MessageEventManager:BaseEventManager<MessageEvent>
    {
        public static MessageEventManager Singleton
        {
            get;
        } = new MessageEventManager();

        private MessageEventManager()
        {
            
        }
    }
}