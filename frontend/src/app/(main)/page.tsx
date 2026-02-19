"use client";

import HeroSection from "@/components/landing-page/hero-section";
import Navbar from "@/components/landing-page/navbar";
import UserAvatar from "@/components/user-avatar";

const Page = () => {
  return (
    <div className="min-h-screen">
      <Navbar />
      <HeroSection />
    </div>
  );
};

export default Page;
