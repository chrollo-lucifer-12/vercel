import { ArrowRightIcon } from "@phosphor-icons/react";
import { Button } from "../ui/button";
import Image from "next/image";

const HeroSection = () => {
  return (
    <section className="relative bg-[url('/hero-section-doodles.png')] bg-contain bg-center bg-no-repeat overflow-hidden pt-32 pb-20 md:pt-40 md:pb-32">
      <div className="container relative mx-auto px-6 text-center">
        <h1 className="mx-auto max-w-3xl font-display text-5xl font-bold leading-[1.1] tracking-tight md:text-7xl">
          Deploy your frontend{" "}
          <span className="text-gradient inline-flex flex-col items-center">
            <p className="leading-none">in seconds</p>

            <Image
              src="/title-highlight-2.png"
              alt="highlight"
              width={320}
              height={12}
              className="-mt-1"
            />
          </span>
        </h1>

        <p className="mx-auto mt-6 max-w-xl text-lg leading-relaxed text-muted-foreground">
          Push your code. We handle the rest. Instant previews, global CDN, and
          zero-config deployments for every frontend framework.
        </p>

        <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <Button
            size="lg"
            className="bg-accent text-accent-foreground hover:bg-accent/90 px-8 gap-2 text-base"
          >
            Start Deploying <ArrowRightIcon size={16} />
          </Button>
          <Button size="lg" variant="outline" className="px-8 text-base">
            View Demo
          </Button>
        </div>

        <p className="mt-4 text-sm text-muted-foreground">
          No credit card required · Free tier forever
        </p>
      </div>
    </section>
  );
};

export default HeroSection;
