"use client";

import Link from "next/link";
import Image from "next/image";

import { Button } from "../ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ListIcon } from "@phosphor-icons/react";
import { useScroll, useTransform, motion } from "motion/react";

const navLinks = [
  { label: "Features", href: "#features" },
  { label: "How It Works", href: "#how-it-works" },
  { label: "Pricing", href: "#pricing" },
];

const Navbar = () => {
  const { scrollY } = useScroll();

  const width = useTransform(scrollY, [0, 120], ["100%", "90%"]);
  const y = useTransform(scrollY, [0, 120], [0, 16]);
  const borderRadius = useTransform(scrollY, [0, 120], [0, 999]);
  const shadow = useTransform(
    scrollY,
    [0, 120],
    ["0px 0px 0px rgba(0,0,0,0)", "0px 10px 30px rgba(0,0,0,0.08)"],
  );
  const blur = useTransform(scrollY, [0, 120], ["blur(0px)", "blur(16px)"]);

  return (
    <div className="fixed top-0 left-0 right-0 z-50 flex justify-center">
      <motion.nav
        style={{
          width,
          y,
          borderRadius,
          boxShadow: shadow,
          backdropFilter: blur,
        }}
        className="flex h-16 items-center justify-between px-6 bg-background/70"
      >
        <Link href="/" className="flex items-center">
          <Image src="/logo.png" alt="logo" width={100} height={50} />
        </Link>

        <div className="hidden md:flex items-center gap-8">
          {navLinks.map((link) => (
            <Link
              key={link.label}
              href={link.href}
              className="text-sm font-medium text-muted-foreground hover:text-foreground transition"
            >
              {link.label}
            </Link>
          ))}
        </div>

        <div className="hidden md:flex items-center gap-3">
          <Button variant="ghost" size="sm" asChild>
            <Link href="/signin">Log in</Link>
          </Button>

          <Button
            size="sm"
            className="bg-accent text-accent-foreground hover:bg-accent/90"
            asChild
          >
            <Link href="/dashboard">Get Started</Link>
          </Button>
        </div>

        <div className="md:hidden">
          <Sheet>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon">
                <ListIcon />
              </Button>
            </SheetTrigger>

            <SheetContent className="w-75 px-6">
              <div className="mt-8 flex flex-col gap-6">
                {navLinks.map((link) => (
                  <Link
                    key={link.label}
                    href={link.href}
                    className="text-base font-medium text-muted-foreground hover:text-foreground"
                  >
                    {link.label}
                  </Link>
                ))}
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </motion.nav>
    </div>
  );
};

export default Navbar;
