"use client";

import Feature1 from "@/components/landing-page/feature-1";
import Feature2 from "@/components/landing-page/feature-2";
import Feature3 from "@/components/landing-page/feature-3";
import Footer from "@/components/landing-page/footer";
import HeroSection from "@/components/landing-page/hero-section";
import Navbar from "@/components/landing-page/navbar";
import QnaSection from "@/components/landing-page/qna-section";
import TutorialVideo from "@/components/landing-page/tutorial-video";

const Page = () => {
  return (
    <div className="min-h-screen">
      <Navbar />
      <HeroSection />
      <TutorialVideo />
      <Feature1 />
      <Feature2 />
      <Feature3 />
      <QnaSection />
      <Footer />
    </div>
  );
};

export default Page;
