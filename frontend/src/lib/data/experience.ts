export type ExperienceItem = {
  key: string;
  label: string;
  title: string;
  org: string;
  orgUrl?: string;
  period: string;
  location: string;
  paragraphs?: string[];
  bullets?: string[];
  image: string; // /assets/pX.jpg
  imageLink?: string;
};

export const EXPERIENCE: ExperienceItem[] = [
  {
    key: "bookbeo",
    label: "Full Stack Mobile Developer",
    title: "Full Stack Mobile Developer",
    org: "bookBeo",
    orgUrl: "https://www.bookbeo.com/",
    period: "Mar 2023 – Present (Internship)",
    location: "Rennes, Brittany, France · On-site",
    paragraphs: [
      "Internship focused on building and shipping a mobile application with a full-stack mindset: clean UI, reliable integration, and performance-aware implementation.",
      "Explored and integrated AWS services to improve scalability, delivery workflows, and operational robustness.",
    ],
    image: "/assets/p7.jpg",
  },
  {
    key: "stm",
    label: "Security Project Team Member",
    title: "Security Project Team Member",
    org: "STMicroelectronics",
    period: "Sep 2022 – Mar 2023",
    location: "Le Mans, Pays de la Loire, France · On-site",
    paragraphs: [
      "Study project on the PSA Cryptography API standard and its specifications, with focus on embedded constraints.",
      "Worked close to the ARM Cortex-M ecosystem; implemented secure primitives and tooling in C; versioned and collaborated with Git/Linux tooling.",
    ],
    image: "/assets/p6.jpg",
  },
  {
    key: "falcon",
    label: "Mobile Developer Trainee",
    title: "Mobile Developer Trainee",
    org: "Falcon Inspection & Services",
    period: "Jun 2022 – Aug 2022",
    location: "Tunis, Tunisia · On-site",
    paragraphs: [
      "Built a mobile application to verify and record data from four different tests: Resistance Test, Thermal Stabilization Test, Air Presence Test, and Waterproof Test.",
      "Each test included its own attributes and mathematical equations with validated data entry flows.",
    ],
    bullets: [
      "Front-end: Flutter",
      "Back-end: Node.js",
      "Practices: Docker, Git, automation, DevOps-oriented workflows",
    ],
    image: "/assets/p5.jpg",
  },
  {
    key: "dall-2021",
    label: "Monitoring & Storage of Meteorological Parameters",
    title: "Monitoring and Storage of Meteorological Parameters of a Given Area",
    org: "DigiArtLivingLab (DALL)",
    period: "9 Jul 2021 – 23 Aug 2021",
    location: "Tunis, Tunisia · On-site",
    bullets: [
      "Built automation + computation utilities using Java.",
      "Implemented server-side logic to automate storage of collected data.",
      "Used OpenWeather APIs; parsed JSON payloads reliably.",
      "Created a user-friendly interface with JavaFX and Swing.",
    ],
    image: "/assets/p1.jpg",
  },
  {
    key: "sagemcom",
    label: "Validate Electronic Cards",
    title: "Verify the Existence of Electronic Components and Validate Cards",
    org: "SAGEMCOM",
    period: "Feb 2020 – May 2020",
    location: "Tunis, Tunisia · On-site",
    bullets: [
      "Used LabVIEW libraries to build a component selection + detection workflow.",
      "Trained a CNN model for multi-component detection.",
      "Converted the model to TensorFlow Lite and ran it on Raspberry Pi.",
      "Designed a camera-holder part in SolidWorks.",
    ],
    image: "/assets/p2.jpg",
    imageLink: "https://drive.google.com/file/d/1HEhe-tPGshCTeTujwacxx_Wc-qHJIF3j/view",
  },
  {
    key: "dall-2019",
    label: "Arduino Stations + Android Apps",
    title: "Humidity/Temperature Station + Scent Diffuser + Android Apps",
    org: "DigiArtLivingLab (DALL)",
    period: "Jan 2019 – Feb 2019",
    location: "Tunis, Tunisia · On-site",
    bullets: [
      "Built Arduino-based devices; implemented Bluetooth data transmission.",
      "Designed and assembled both the station and diffuser with available tools.",
      "Created both apps in MIT App Inventor under time constraints.",
      "Linked apps to Firebase and stored data in real time.",
    ],
    image: "/assets/p3.jpg",
  },
  {
    key: "dall-2018",
    label: "Educational Direction Game",
    title: "Educational Game to Teach Directions (Unity)",
    org: "DigiArtLivingLab (DALL)",
    period: "Jan 2018 – Feb 2018",
    location: "Nabeul, Tunisia · On-site",
    bullets: [
      "Worked with Unity and asset import pipelines.",
      "Implemented gameplay components (C/C++ exposure during the project).",
      "Collaborated in a team setting; learned iteration and delivery discipline.",
      "Designed a maze level and tuned audiovisual cues to keep engagement high.",
    ],
    image: "/assets/p4.jpg",
  },
];
