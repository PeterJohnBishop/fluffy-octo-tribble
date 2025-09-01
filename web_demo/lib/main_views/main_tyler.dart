import 'package:flutter/material.dart';

class TylerView extends StatelessWidget {
  const TylerView({super.key});

  final Map<String, dynamic> realtorData = const {
    "assetPath": "/tyler.jpg",
    "name": "Tyler Joyce",
    "title": "Associate",
    "License": "#FA.100106696",
    "experience": "1 Year",
    "languages": "English",
    "specialties": "Military, First Time Home Buyers, Investments, Historic Homes",
    "bio": ["Coming soon..."]
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
