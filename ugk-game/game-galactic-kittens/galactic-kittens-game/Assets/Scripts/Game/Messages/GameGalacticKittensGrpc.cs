// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: game_galactic_kittens.proto
// </auto-generated>
#pragma warning disable 0414, 1591
#region Designer generated code

using grpc = global::Grpc.Core;

/// <summary>
///GalacticKittens Match 服务
/// </summary>
public static partial class GalacticKittensMatchService
{
  static readonly string __ServiceName = "GalacticKittensMatchService";

  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static void __Helper_SerializeMessage(global::Google.Protobuf.IMessage message, grpc::SerializationContext context)
  {
    #if !GRPC_DISABLE_PROTOBUF_BUFFER_SERIALIZATION
    if (message is global::Google.Protobuf.IBufferMessage)
    {
      context.SetPayloadLength(message.CalculateSize());
      global::Google.Protobuf.MessageExtensions.WriteTo(message, context.GetBufferWriter());
      context.Complete();
      return;
    }
    #endif
    context.Complete(global::Google.Protobuf.MessageExtensions.ToByteArray(message));
  }

  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static class __Helper_MessageCache<T>
  {
    public static readonly bool IsBufferMessage = global::System.Reflection.IntrospectionExtensions.GetTypeInfo(typeof(global::Google.Protobuf.IBufferMessage)).IsAssignableFrom(typeof(T));
  }

  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static T __Helper_DeserializeMessage<T>(grpc::DeserializationContext context, global::Google.Protobuf.MessageParser<T> parser) where T : global::Google.Protobuf.IMessage<T>
  {
    #if !GRPC_DISABLE_PROTOBUF_BUFFER_SERIALIZATION
    if (__Helper_MessageCache<T>.IsBufferMessage)
    {
      return parser.ParseFrom(context.PayloadAsReadOnlySequence());
    }
    #endif
    return parser.ParseFrom(context.PayloadAsNewBuffer());
  }

  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static readonly grpc::Marshaller<global::GalacticKittensPlayerServerListRequest> __Marshaller_GalacticKittensPlayerServerListRequest = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::GalacticKittensPlayerServerListRequest.Parser));
  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static readonly grpc::Marshaller<global::GalacticKittensPlayerServerListResponse> __Marshaller_GalacticKittensPlayerServerListResponse = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::GalacticKittensPlayerServerListResponse.Parser));
  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static readonly grpc::Marshaller<global::GalacticKittensGameFinishRequest> __Marshaller_GalacticKittensGameFinishRequest = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::GalacticKittensGameFinishRequest.Parser));
  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static readonly grpc::Marshaller<global::GalacticKittensGameFinishResponse> __Marshaller_GalacticKittensGameFinishResponse = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::GalacticKittensGameFinishResponse.Parser));

  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static readonly grpc::Method<global::GalacticKittensPlayerServerListRequest, global::GalacticKittensPlayerServerListResponse> __Method_playerServerList = new grpc::Method<global::GalacticKittensPlayerServerListRequest, global::GalacticKittensPlayerServerListResponse>(
      grpc::MethodType.Unary,
      __ServiceName,
      "playerServerList",
      __Marshaller_GalacticKittensPlayerServerListRequest,
      __Marshaller_GalacticKittensPlayerServerListResponse);

  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  static readonly grpc::Method<global::GalacticKittensGameFinishRequest, global::GalacticKittensGameFinishResponse> __Method_gameFinish = new grpc::Method<global::GalacticKittensGameFinishRequest, global::GalacticKittensGameFinishResponse>(
      grpc::MethodType.Unary,
      __ServiceName,
      "gameFinish",
      __Marshaller_GalacticKittensGameFinishRequest,
      __Marshaller_GalacticKittensGameFinishResponse);

  /// <summary>Service descriptor</summary>
  public static global::Google.Protobuf.Reflection.ServiceDescriptor Descriptor
  {
    get { return global::GameGalacticKittensReflection.Descriptor.Services[0]; }
  }

  /// <summary>Base class for server-side implementations of GalacticKittensMatchService</summary>
  [grpc::BindServiceMethod(typeof(GalacticKittensMatchService), "BindService")]
  public abstract partial class GalacticKittensMatchServiceBase
  {
    /// <summary>
    ///玩家服务器列表
    /// </summary>
    /// <param name="request">The request received from the client.</param>
    /// <param name="context">The context of the server-side call handler being invoked.</param>
    /// <returns>The response to send back to the client (wrapped by a task).</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual global::System.Threading.Tasks.Task<global::GalacticKittensPlayerServerListResponse> playerServerList(global::GalacticKittensPlayerServerListRequest request, grpc::ServerCallContext context)
    {
      throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
    }

    /// <summary>
    /// 游戏完成
    /// </summary>
    /// <param name="request">The request received from the client.</param>
    /// <param name="context">The context of the server-side call handler being invoked.</param>
    /// <returns>The response to send back to the client (wrapped by a task).</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual global::System.Threading.Tasks.Task<global::GalacticKittensGameFinishResponse> gameFinish(global::GalacticKittensGameFinishRequest request, grpc::ServerCallContext context)
    {
      throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
    }

  }

  /// <summary>Client for GalacticKittensMatchService</summary>
  public partial class GalacticKittensMatchServiceClient : grpc::ClientBase<GalacticKittensMatchServiceClient>
  {
    /// <summary>Creates a new client for GalacticKittensMatchService</summary>
    /// <param name="channel">The channel to use to make remote calls.</param>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public GalacticKittensMatchServiceClient(grpc::ChannelBase channel) : base(channel)
    {
    }
    /// <summary>Creates a new client for GalacticKittensMatchService that uses a custom <c>CallInvoker</c>.</summary>
    /// <param name="callInvoker">The callInvoker to use to make remote calls.</param>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public GalacticKittensMatchServiceClient(grpc::CallInvoker callInvoker) : base(callInvoker)
    {
    }
    /// <summary>Protected parameterless constructor to allow creation of test doubles.</summary>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    protected GalacticKittensMatchServiceClient() : base()
    {
    }
    /// <summary>Protected constructor to allow creation of configured clients.</summary>
    /// <param name="configuration">The client configuration.</param>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    protected GalacticKittensMatchServiceClient(ClientBaseConfiguration configuration) : base(configuration)
    {
    }

    /// <summary>
    ///玩家服务器列表
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="headers">The initial metadata to send with the call. This parameter is optional.</param>
    /// <param name="deadline">An optional deadline for the call. The call will be cancelled if deadline is hit.</param>
    /// <param name="cancellationToken">An optional token for canceling the call.</param>
    /// <returns>The response received from the server.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual global::GalacticKittensPlayerServerListResponse playerServerList(global::GalacticKittensPlayerServerListRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
    {
      return playerServerList(request, new grpc::CallOptions(headers, deadline, cancellationToken));
    }
    /// <summary>
    ///玩家服务器列表
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="options">The options for the call.</param>
    /// <returns>The response received from the server.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual global::GalacticKittensPlayerServerListResponse playerServerList(global::GalacticKittensPlayerServerListRequest request, grpc::CallOptions options)
    {
      return CallInvoker.BlockingUnaryCall(__Method_playerServerList, null, options, request);
    }
    /// <summary>
    ///玩家服务器列表
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="headers">The initial metadata to send with the call. This parameter is optional.</param>
    /// <param name="deadline">An optional deadline for the call. The call will be cancelled if deadline is hit.</param>
    /// <param name="cancellationToken">An optional token for canceling the call.</param>
    /// <returns>The call object.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual grpc::AsyncUnaryCall<global::GalacticKittensPlayerServerListResponse> playerServerListAsync(global::GalacticKittensPlayerServerListRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
    {
      return playerServerListAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
    }
    /// <summary>
    ///玩家服务器列表
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="options">The options for the call.</param>
    /// <returns>The call object.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual grpc::AsyncUnaryCall<global::GalacticKittensPlayerServerListResponse> playerServerListAsync(global::GalacticKittensPlayerServerListRequest request, grpc::CallOptions options)
    {
      return CallInvoker.AsyncUnaryCall(__Method_playerServerList, null, options, request);
    }
    /// <summary>
    /// 游戏完成
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="headers">The initial metadata to send with the call. This parameter is optional.</param>
    /// <param name="deadline">An optional deadline for the call. The call will be cancelled if deadline is hit.</param>
    /// <param name="cancellationToken">An optional token for canceling the call.</param>
    /// <returns>The response received from the server.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual global::GalacticKittensGameFinishResponse gameFinish(global::GalacticKittensGameFinishRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
    {
      return gameFinish(request, new grpc::CallOptions(headers, deadline, cancellationToken));
    }
    /// <summary>
    /// 游戏完成
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="options">The options for the call.</param>
    /// <returns>The response received from the server.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual global::GalacticKittensGameFinishResponse gameFinish(global::GalacticKittensGameFinishRequest request, grpc::CallOptions options)
    {
      return CallInvoker.BlockingUnaryCall(__Method_gameFinish, null, options, request);
    }
    /// <summary>
    /// 游戏完成
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="headers">The initial metadata to send with the call. This parameter is optional.</param>
    /// <param name="deadline">An optional deadline for the call. The call will be cancelled if deadline is hit.</param>
    /// <param name="cancellationToken">An optional token for canceling the call.</param>
    /// <returns>The call object.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual grpc::AsyncUnaryCall<global::GalacticKittensGameFinishResponse> gameFinishAsync(global::GalacticKittensGameFinishRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
    {
      return gameFinishAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
    }
    /// <summary>
    /// 游戏完成
    /// </summary>
    /// <param name="request">The request to send to the server.</param>
    /// <param name="options">The options for the call.</param>
    /// <returns>The call object.</returns>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public virtual grpc::AsyncUnaryCall<global::GalacticKittensGameFinishResponse> gameFinishAsync(global::GalacticKittensGameFinishRequest request, grpc::CallOptions options)
    {
      return CallInvoker.AsyncUnaryCall(__Method_gameFinish, null, options, request);
    }
    /// <summary>Creates a new instance of client from given <c>ClientBaseConfiguration</c>.</summary>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    protected override GalacticKittensMatchServiceClient NewInstance(ClientBaseConfiguration configuration)
    {
      return new GalacticKittensMatchServiceClient(configuration);
    }
  }

  /// <summary>Creates service definition that can be registered with a server</summary>
  /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  public static grpc::ServerServiceDefinition BindService(GalacticKittensMatchServiceBase serviceImpl)
  {
    return grpc::ServerServiceDefinition.CreateBuilder()
        .AddMethod(__Method_playerServerList, serviceImpl.playerServerList)
        .AddMethod(__Method_gameFinish, serviceImpl.gameFinish).Build();
  }

  /// <summary>Register service method with a service binder with or without implementation. Useful when customizing the  service binding logic.
  /// Note: this method is part of an experimental API that can change or be removed without any prior notice.</summary>
  /// <param name="serviceBinder">Service methods will be bound by calling <c>AddMethod</c> on this object.</param>
  /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
  [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
  public static void BindService(grpc::ServiceBinderBase serviceBinder, GalacticKittensMatchServiceBase serviceImpl)
  {
    serviceBinder.AddMethod(__Method_playerServerList, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::GalacticKittensPlayerServerListRequest, global::GalacticKittensPlayerServerListResponse>(serviceImpl.playerServerList));
    serviceBinder.AddMethod(__Method_gameFinish, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::GalacticKittensGameFinishRequest, global::GalacticKittensGameFinishResponse>(serviceImpl.gameFinish));
  }

}
#endregion
