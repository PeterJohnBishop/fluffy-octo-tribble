import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class MainNavView extends StatelessWidget {
  const MainNavView({
    super.key,
    required this.height,
    required this.width,
    required this.selectedIndex,
    required this.onIndexChanged,
    required this.ws,
  });

  final double height;
  final double width;
  final int selectedIndex;
  final ValueChanged<int> onIndexChanged;
  final ws;

  @override
  Widget build(BuildContext context) {
    final navOptions = ["Home", "About", "Listings", "Buyers", "Sellers", "Contact"];
    return SizedBox(
      height: height * 0.05,
      width: width,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: List.generate(navOptions.length, (index) {
          final isSelected = index == selectedIndex;
          return TextButton(
            onPressed: () {
              onIndexChanged(index);
              ws.send({
                "event": "web",
                "message": "${navOptions[index]} viewed!",
                "timestamp": DateTime.now().millisecondsSinceEpoch, 
              });
              },
            child: Text(
              navOptions[index],
              style: GoogleFonts.alumniSansPinstripe().copyWith(
                fontSize: 20,
                fontWeight: FontWeight.w600,
                color: Colors.black,
                decoration: isSelected ? TextDecoration.underline : TextDecoration.none,
              ),
            ),
          );
        }),
      ),
    );
  }
}
