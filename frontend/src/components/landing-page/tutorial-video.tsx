"use client";

import { motion } from "motion/react";

const TutorialVideo = () => {
  return (
    <div className="flex items-center justify-center mb-10">
      <motion.div
        initial={{ opacity: 0, y: 40, scale: 0.95 }}
        whileInView={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
        viewport={{ once: true }}
        whileHover={{ y: -6 }}
        className="border border-gray-400 w-full sm:w-[60%] rounded-md shadow-md bg-white"
      >
        <motion.div
          className="h-10 flex gap-x-2 items-center p-2"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4 }}
        >
          <motion.div
            className="rounded-full w-[10px] h-[10px] bg-gray-400"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.5, type: "spring" }}
          />
          <motion.div
            className="rounded-full w-[10px] h-[10px] bg-gray-400"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.6, type: "spring" }}
          />
          <motion.div
            className="rounded-full w-[10px] h-[10px] bg-gray-400"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.7, type: "spring" }}
          />
        </motion.div>

        <motion.div
          className="border border-gray-400 overflow-hidden"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.8, duration: 0.5 }}
        >
          <video playsInline muted autoPlay loop width="100%">
            <source src="/intro.mp4" type="video/mp4" />
          </video>
        </motion.div>
      </motion.div>
    </div>
  );
};

export default TutorialVideo;
