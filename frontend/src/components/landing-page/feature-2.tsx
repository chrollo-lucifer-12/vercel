"use client";

import Image from "next/image";
import { useEffect, useState } from "react";

const images = ["inputs-1.png", "inputs-4.png", "inputs-6.png"];

const Feature2 = () => {
  const [currImage, setCurrImage] = useState(0);

  useEffect(() => {
    const intervalId = setInterval(() => {
      setCurrImage((prevIndex) => (prevIndex + 1) % images.length);
    }, 2000);

    return () => clearInterval(intervalId);
  }, []);
  return (
    <section className="mt-32 flex flex-col items-center  gap-y-10 px-4 sm:px-10 md:px-20 lg:px-[120px]">
      <div
        className={
          "sm:w-[600px] md:w-[700px] w-[300px] flex justify-between items-center"
        }
      >
        <div>
          <h1 className="text-black font-extrabold sm:text-3xl text-lg">
            Powerful deployments. Zero complexity.
          </h1>
          <p className="mt-2 text-gray-700  text-[12px] max-w-[250px] sm:text-[16px] sm:max-w-[1200px]">
            Everything you need to ship modern frontends — automatic builds,
            preview environments, global CDN delivery, and instant rollbacks.
            All without touching infrastructure.
          </p>
        </div>
        <Image
          src={"/smart.png"}
          alt={"smart"}
          width={250}
          height={250}
          className={"hidden sm:block"}
        />
      </div>
      <div
        className={
          "rounded-2xl sm:w-[600px] md:w-[700px] w-[300px] h-auto border border-gray-6 p-10 text-black"
        }
      >
        <Image
          src={"/icon2.svg"}
          alt={"icon1"}
          width={20}
          height={20}
          className={"text-pink-500"}
        />
        <h1 className={"mt-2 font-extrabold text-black"}>
          Deploy from Git in seconds
        </h1>
        <p className={"text-xs text-black mt-1"}>
          Connect your GitHub repository and every push triggers a
          production-ready build. Get preview URLs for pull requests, automatic
          framework detection, optimized asset delivery, and blazing-fast edge
          performance out of the box.
        </p>
        <div className={" border border-gray-4 mt-10 rounded-md shadow-md"}>
          <div className={"h-10 flex gap-x-2 items-center p-2"}>
            <div className={"rounded-full w-[10px] h-[10px] bg-gray-6"} />
            <div className={"rounded-full w-[10px] h-[10px] bg-gray-6"} />
            <div className={"rounded-full w-[10px] h-[10px] bg-gray-6"} />
          </div>
          <div className={"border border-gray-4"}>
            <Image
              src={`/${images[currImage]}`}
              alt={`form image ${currImage}`}
              width={500}
              height={500}
              className={"transition-opacity duration-700 ease-in-out"}
            />
          </div>
        </div>
      </div>
    </section>
  );
};

export default Feature2;
