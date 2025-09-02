import 'package:flutter/material.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:url_launcher/url_launcher.dart';

class MainCarouselView extends StatelessWidget {
  const MainCarouselView({super.key});

  final List<Map<String, String>> carouselItems = const [
    {
      "imageUrl": "https://images.unsplash.com/photo-1505691938895-1758d7feb511",
      "title": "Spacious Living Room",
      "subtitle": "Open concept with natural lighting",
      "linkUrl": "https://example.com/property1"
    },
    {
      "imageUrl": "https://images.unsplash.com/photo-1568605114967-8130f3a36994",
      "title": "Warm and Modern Exterior",
      "subtitle": "Landscaped yard with a patio",
      "linkUrl": "https://example.com/property2"
    },
    {
      "imageUrl": "https://images.unsplash.com/photo-1507089947368-19c1da9775ae",
      "title": "Modern Kitchen",
      "subtitle": "High-end appliances and finishes",
      "linkUrl": "https://example.com/property3"
    },
    {
      "imageUrl": "https://images.unsplash.com/photo-1600585152220-90363fe7e115",
      "title": "Basement Kitchen",
      "subtitle": "High-end appliances and finishes",
      "linkUrl": "https://example.com/property4"
    },
    {
      "imageUrl": "https://images.unsplash.com/photo-1494526585095-c41746248156",
      "title": "Luxury Exterior",
      "subtitle": "Elegant design with curb appeal",
      "linkUrl": "https://example.com/property5"
    },
  ];

  void _launchURL(String url) async {
    final uri = Uri.parse(url);
    if (await canLaunchUrl(uri)) {
      await launchUrl(uri);
    } else {
      throw 'Could not launch $url';
    }
  }

  @override
  Widget build(BuildContext context) {
    return CarouselSlider(
        options: CarouselOptions(
          height: double.infinity,
          viewportFraction: 1.0,
          enlargeCenterPage: false,
          autoPlay: true,
          autoPlayInterval: const Duration(seconds: 5),
          autoPlayAnimationDuration: const Duration(milliseconds: 800),
          enableInfiniteScroll: true,
        ),
        items: carouselItems.map((item) {
          return Builder(
            builder: (BuildContext context) {
              return Stack(
                fit: StackFit.expand,
                children: [
                  Image.network(
                    item["imageUrl"]!,
                    fit: BoxFit.cover,
                  ),
                  Container(
                    color: Colors.black.withOpacity(0.4),
                  ),
                  Positioned(
                    left: 20,
                    bottom: 40,
                    right: 20,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          item["title"]!,
                          style: const TextStyle(
                            fontSize: 28,
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          item["subtitle"]!,
                          style: const TextStyle(
                            fontSize: 18,
                            color: Colors.white70,
                          ),
                        ),
                        const SizedBox(height: 16),
                        ElevatedButton(
                          onPressed: () => _launchURL(item["linkUrl"]!),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.white,
                            foregroundColor: Colors.black,
                          ),
                          child: const Text("View Property"),
                        ),
                      ],
                    ),
                  ),
                ],
              );
            },
          );
        }).toList(),
    );
  }
}

