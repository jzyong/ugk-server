using System.Collections;

namespace Game.Room
{
    public interface IDamagable
    {
        void Hit(int damage);

        IEnumerator HitEffect();
    }
}