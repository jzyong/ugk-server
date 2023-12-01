using System.Collections;
using Game.Manager;
using UnityEngine;

namespace Game.Room
{
    /// <summary>
    /// 防御护盾 
    /// </summary>
    public class DefenseMatrix : MonoBehaviour, IDamagable
    {
        public bool isShieldActive { get; private set; } = false;


        private void Start()
        {
        }

        public void Hit(int damage)
        {
            //TODO 通知客户端关闭护盾
            isShieldActive = false;
        }

        public void TurnOnShield()
        {
            isShieldActive = true;
        }

        private void OnTriggerEnter2D(Collider2D collider)
        {
            if (collider.TryGetComponent(out IDamagable damagable))
            {
                damagable.Hit(1);
                //TODO 通知客户端关闭护盾
                isShieldActive = false;

                GalacticKittensUseShieldResponse response = new GalacticKittensUseShieldResponse()
                {
                    //TODO 获取所属的飞船ID
                    State = 1
                };
                PlayerManager.Singleton.BroadcastMsg(MID.GalacticKittensUseShieldRes, response);
            }
        }


        IEnumerator IDamagable.HitEffect()
        {
            return null;
        }
    }
}