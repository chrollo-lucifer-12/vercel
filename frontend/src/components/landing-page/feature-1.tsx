"use client";

import Image from "next/image";
import { motion } from "motion/react";

const container = {
  hidden: {},
  show: {
    transition: {
      staggerChildren: 0.15,
    },
  },
};

const fadeUp = {
  hidden: { opacity: 0, y: 40 },
  show: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.7,
      ease: [0.22, 1, 0.36, 1] as const,
    },
  },
};

const cardAnimation = {
  hidden: { opacity: 0, scale: 0.9, y: 60 },
  show: {
    opacity: 1,
    scale: 1,
    y: 0,
    transition: {
      duration: 0.8,
      ease: [0.22, 1, 0.36, 1] as const,
    },
  },
};

const Feature1 = () => {
  return (
    <motion.section
      className="mt-32 flex flex-col items-center gap-y-10 sm:px-10 md:px-20 lg:px-[120px] mb-10"
      variants={container}
      initial="hidden"
      whileInView="show"
      viewport={{ once: true, amount: 0.2 }}
    >
      <motion.div
        className="sm:w-[600px] md:w-[700px] w-[300px]"
        variants={fadeUp}
      >
        <motion.h1
          variants={fadeUp}
          className="text-black font-extrabold text-lg sm:text-3xl"
        >
          Deploy your frontend in seconds
        </motion.h1>

        <motion.p
          variants={fadeUp}
          className="mt-2 text-gray-700 text-[12px] max-w-[250px] sm:text-[12px] md:text-[16px] sm:max-w-[650px] md:max-w-[1200px]"
        >
          Push your code and go live instantly. Built for modern frameworks like
          Next.js, React, and Vue — with zero configuration required.
        </motion.p>
      </motion.div>

      <motion.div
        variants={cardAnimation}
        whileHover={{ scale: 1.03 }}
        className="rounded-2xl sm:w-[600px] md:w-[700px] w-[300px] h-auto border p-10 text-black shadow-[0_0_0_2px_rgb(248,28,229),0_0_0_4px_rgba(248,28,229,0.36)]"
      >
        <motion.div
          animate={{ y: [0, -8, 0] }}
          transition={{
            duration: 3,
            repeat: Infinity,
            ease: "easeInOut",
          }}
        >
          <Image
            src="/dive-in.png"
            alt="dive in"
            width={200}
            height={100}
            className="w-full"
          />
        </motion.div>

        <motion.p
          variants={fadeUp}
          className="text-[18px] font-[1000] bg-gradient-to-r from-[#8A46FF] to-[#F81CE0] bg-clip-text text-transparent mt-4"
        >
          Unlimited deployments. Zero friction.
        </motion.p>

        <motion.p variants={fadeUp} className="max-w-[620px] mt-1 text-sm">
          Deploy directly from GitHub with automatic builds, global CDN
          delivery, preview URLs for every commit, and instant rollbacks. Focus
          on building — we handle the infrastructure.
        </motion.p>
      </motion.div>
    </motion.section>
  );
};

export default Feature1;
