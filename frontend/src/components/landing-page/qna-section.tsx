"use client";

import { ArrowRightIcon, CaretDownIcon } from "@phosphor-icons/react";
import { useState } from "react";

const questions = [
  {
    question: "Is your platform really free?",
    answer:
      "Yes! You can deploy unlimited personal projects for free within our fair usage limits. Connect your Git repository and get automatic builds, preview URLs, and global CDN delivery — all without entering a credit card. Upgrade only when you need advanced features or higher limits.",
  },
  {
    question: "Are deployments secure?",
    answer:
      "Absolutely. All deployments are served over HTTPS by default with automatic SSL certificates. Your builds run in isolated environments, and your source code remains securely connected through your Git provider. We also ensure encrypted data transfer and secure infrastructure best practices.",
  },
  {
    question: "How does this compare to other deployment platforms?",
    answer:
      "Unlike traditional hosting providers that require manual configuration, we automatically detect your framework, build your app, and deploy it globally in seconds. Every push creates a production-ready deployment and preview link — no DevOps setup required.",
  },
  {
    question: "How can I get started?",
    answer:
      "Just sign up, connect your GitHub repository, and push your code. Your project will automatically build and deploy in seconds. No configuration needed — your frontend goes live instantly.",
  },
];

const QnaSection = () => {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggleAnswer = (index: number) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  return (
    <section className="mt-32 flex flex-col items-center gap-y-10 px-4 sm:px-10 md:px-20 lg:px-[120px] mb-10">
      <div className="sm:w-[600px] md:w-[700px] w-[300px]">
        <h1 className="text-black font-extrabold sm:text-3xl text-lg">
          Questions & answers
        </h1>
      </div>
      <div className="sm:w-[600px] md:w-[700px] w-[300px] mt-2">
        {questions.map((q, i) => (
          <div key={i} className="w-full border-t border-gray-300">
            <button
              onClick={() => toggleAnswer(i)}
              className="w-full flex justify-between items-center text-left p-4 cursor-pointer"
            >
              <p className="font-bold text-black text-sm">{q.question}</p>
              {openIndex === i ? (
                <CaretDownIcon className="text-black transition-transform rotate-180" />
              ) : (
                <ArrowRightIcon className="text-black" />
              )}
            </button>
            {openIndex === i && (
              <div className="p-4 text-sm text-gray-700 border-t border-gray-200 bg-gray-50">
                {q.answer}
              </div>
            )}
          </div>
        ))}
      </div>
    </section>
  );
};

export default QnaSection;
