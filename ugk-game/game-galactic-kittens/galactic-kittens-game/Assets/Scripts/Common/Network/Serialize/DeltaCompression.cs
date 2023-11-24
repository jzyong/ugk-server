// manual delta compression for some types.
//    varint(b-a)
// Mirror can't use Mirror II's bit-tree delta compression.

using System.Runtime.CompilerServices;
using Common.Tools;

namespace Common.Network.Serialize
{
    /// <summary>
    /// [Mirror](https://github.com/MirrorNetworking/Mirror) 参考`Compression.cs`,`DeltaCompression.cs`,`Vector3Long.cs`  
    ///实现流程如下：  
    ///1. 对象坐标 Vector3 a=(103.1,35.2,221.2) 变为 (105.5,40,223) 
    ///2. 将坐标进行精度保留变为整数坐标 (1031,352,2212)(1055,400,2230)
    ///3. 计算变化值(24,48,18)
    ///4. 字节宽度压缩，因为三个坐标都小于240，因此只需要三个Byte就能传输
    ///5. 传输字节从3个float变为3个byte，共减少3*4-3*1=9Byte
    ///6. 接收端根据历史信息还原出真实的坐标
    /// </summary>
    public static class DeltaCompression
    {
        // delta (usually small), then zigzag varint to support +- changes
        // parameter order: (last, current) makes most sense (Q3 does this too).
        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static void Compress(NetworkWriter writer, long last, long current) =>
            Compression.CompressVarInt(writer, current - last);

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static long Decompress(NetworkReader reader, long last) =>
            last + Compression.DecompressVarInt(reader);

        // delta (usually small), then zigzag varint to support +- changes
        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static void Compress(NetworkWriter writer, Vector3Long last, Vector3Long current)
        {
            Compress(writer, last.x, current.x);
            Compress(writer, last.y, current.y);
            Compress(writer, last.z, current.z);
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static Vector3Long Decompress(NetworkReader reader, Vector3Long last)
        {
            long x = Decompress(reader, last.x);
            long y = Decompress(reader, last.y);
            long z = Decompress(reader, last.z);
            return new Vector3Long(x, y, z);
        }
    }
}