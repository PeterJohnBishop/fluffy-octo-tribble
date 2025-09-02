import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';

class WebSocketService {
  static final WebSocketService _instance = WebSocketService._internal();
  factory WebSocketService() => _instance;

  WebSocketService._internal();

  late WebSocketChannel _channel;

  void connect(String url) {
    _channel = WebSocketChannel.connect(Uri.parse(url));
      print("Connected to $url");
    _channel.stream.listen((event) {
      print("Data Received: $event");
    }, onError: (error) {
      print("Error: $error");
    }, onDone: () {
      print("Disconnected from $url");
    });
  }

  void send(Map<String, dynamic> message) {
    final jsonMessage = jsonEncode(message);
    _channel.sink.add(jsonMessage);
  }

  void disconnect() {
    _channel.sink.close();
  }
}

