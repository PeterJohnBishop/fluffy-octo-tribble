import 'package:flutter/material.dart';

class EddieView extends StatelessWidget {
  const EddieView({super.key});

  final Map<String, dynamic> realtorData = const {
    "assetPath": "/eddie.jpg",
    "name": "Edward Lederman",
    "title": "Broker / Co Owner",
    "License": "#FA100080975",
    "experience": "7 Years",
    "languages": "English",
    "specialties": "Business Opportunities",
    "bio": [
      "Eddie Lederman is a proud Denver local who spent his childhood enjoying the breathtaking scenery of the Rocky Mountains. Prior to embarking on his successful career in real estate, Eddie was an Upper Elementary Montessori teacher, which endowed him with exceptional communication skills to effectively convey client needs and represent them tenaciously."
      "Growing up in the vibrant Park Hill neighborhood near the Museum of Nature and Science and City Park Lake, Eddie currently resides on the west side of Denver. He spends his leisure time training for triathlons with his dogs on Weir Gulch, and he enjoys community gardening at the Lowell Street Garden. Eddie is a people person who embraces diversity and enjoys engaging in meaningful conversations to serve the community.",
      "As a realtor, Eddie finds great satisfaction in empowering his clients through ownership. He prioritizes his clients’ needs and ensures they feel satisfied with their purchase. He is a skilled advocate and a verbal learner who enjoys conversing with his clients at length to ensure that they have all the information they need to make informed decisions.",
      "To stay up-to-date with Denver’s dynamic real estate market, Eddie attends monthly governance meetings for the Denver Metro Realtor Association. He values collaboration and innovation, which is why he consistently surrounds himself with like-minded professionals who are passionate about serving their clients’ best interests. As a former teacher, Eddie is adept at creating game plans for his team to ensure everyone is moving in the same direction. Whether it’s scheduling staging and repair services or making sure the lender and underwriter and insurance agent are meeting deadlines, he will be direct, genuine, and empathetic to his clients’ needs and always puts the person before the profit.",
      "Eddie is available 24/7, even when he’s out hiking or cross-country skiing. Contact him to embark on a successful real estate journey today!"
    ]
  };

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(builder: (context, constraints) {
      final bool isWideScreen = constraints.maxWidth > 800;

      Widget imageWidget = Container(
        width: isWideScreen ? constraints.maxWidth * 0.5 : double.infinity,
        padding: const EdgeInsets.all(8),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: Image.asset(
            realtorData["assetPath"],
            fit: BoxFit.cover,
          ),
        ),
      );

      Widget textWidget = Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            realtorData["name"],
            style: const TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            realtorData["title"],
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.w500,
              color: Colors.grey,
            ),
          ),
          const SizedBox(height: 4),
          Text("License: ${realtorData["License"]}", style: const TextStyle(fontSize: 16)),
          const SizedBox(height: 4),
          Text("Experience: ${realtorData["experience"]}", style: const TextStyle(fontSize: 16)),
          const SizedBox(height: 4),
          Text("Languages: ${realtorData["languages"]}", style: const TextStyle(fontSize: 16)),
          const SizedBox(height: 4),
          Text("Specialties: ${realtorData["specialties"]}", style: const TextStyle(fontSize: 16)),
          const SizedBox(height: 16),
          ...List<Widget>.from(
            (realtorData["bio"] as List<String>).map(
              (paragraph) => Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: Text(
                  paragraph,
                  style: const TextStyle(fontSize: 16),
                ),
              ),
            ),
          ),
        ],
      );

      return Padding(
        padding: const EdgeInsets.all(16.0),
        child: isWideScreen
            ? Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(child: imageWidget),
                  const SizedBox(width: 24),
                  Expanded(child: textWidget),
                ],
              )
            : Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  imageWidget,
                  const SizedBox(height: 16),
                  textWidget,
                ],
              ),
      );
    });
  }
}