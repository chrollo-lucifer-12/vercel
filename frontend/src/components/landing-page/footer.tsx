"use client";

import {
  FacebookLogoIcon,
  LinkedinLogoIcon,
  TwitterLogoIcon,
  YoutubeLogoIcon,
} from "@phosphor-icons/react";
import Image from "next/image";
import Link from "next/link";

const Footer = () => {
  const socialIcons = [
    { icon: <TwitterLogoIcon size={18} />, label: "Twitter" },
    { icon: <LinkedinLogoIcon size={18} />, label: "LinkedIn" },
    { icon: <YoutubeLogoIcon size={18} />, label: "YouTube" },
    { icon: <FacebookLogoIcon size={18} />, label: "Facebook" },
  ];

  const productLinks = [
    { label: "Features", href: "#" },
    { label: "Pricing", href: "#" },
    { label: "Customers", href: "#" },
    { label: "What's new", href: "#" },
    { label: "Roadmap", href: "#" },
    { label: "Feature requests", href: "#" },
    { label: "Templates", href: "#" },
    { label: "Integrations", href: "#" },
    { label: "Status", href: "#" },
  ];

  const helpLinks = [
    { label: "Get started", href: "#" },
    { label: "How-to guides", href: "#" },
    { label: "Help center", href: "#" },
    { label: "Contact support", href: "#" },
  ];

  const companyLinks = [
    { label: "About us", href: "#" },
    { label: "Blog", href: "#" },
    { label: "Careers", href: "#" },
  ];

  const resourceLinks = [
    { label: "Community", href: "#" },
    { label: "Referral program", href: "#" },
    { label: "Fair use policy", href: "#" },
    { label: "GDPR", href: "#" },
    { label: "Terms & Privacy", href: "#" },
  ];

  return (
    <section className="mt-40 mb-24 px-6 sm:px-10 md:px-20 lg:px-[120px]">
      {/* Divider */}
      <div className="w-full border-t border-gray-200" />

      <div className="w-full grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-12 pt-12">
        {/* Logo & Info */}
        <div>
          <Image src="/logo.png" alt="logo" width={40} height={40} />

          <p className="text-sm mt-6 text-gray-600 font-medium">
            Global frontend deployments ⚡
          </p>

          <p className="text-sm text-gray-500 mt-1">© 2026 YourCompany Inc.</p>

          <div className="flex gap-5 mt-6 text-gray-500">
            {socialIcons.map((item, index) => (
              <Link
                key={index}
                href="#"
                aria-label={item.label}
                className="hover:text-black transition-colors"
              >
                {item.icon}
              </Link>
            ))}
          </div>
        </div>

        {/* Product */}
        <div>
          <h3 className="font-semibold text-sm mb-4 text-black">Product</h3>
          <div className="flex flex-col gap-3">
            {productLinks.map((link, index) => (
              <Link
                key={index}
                href={link.href}
                className="text-sm text-gray-600 hover:text-black transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </div>
        </div>

        {/* Help + Company */}
        <div>
          <h3 className="font-semibold text-sm mb-4 text-black">Help</h3>
          <div className="flex flex-col gap-3">
            {helpLinks.map((link, index) => (
              <Link
                key={index}
                href={link.href}
                className="text-sm text-gray-600 hover:text-black transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </div>

          <h3 className="font-semibold text-sm mt-10 mb-4 text-black">
            Company
          </h3>
          <div className="flex flex-col gap-3">
            {companyLinks.map((link, index) => (
              <Link
                key={index}
                href={link.href}
                className="text-sm text-gray-600 hover:text-black transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </div>
        </div>

        {/* Resources */}
        <div>
          <h3 className="font-semibold text-sm mb-4 text-black">Resources</h3>
          <div className="flex flex-col gap-3">
            {resourceLinks.map((link, index) => (
              <Link
                key={index}
                href={link.href}
                className="text-sm text-gray-600 hover:text-black transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
};

export default Footer;
