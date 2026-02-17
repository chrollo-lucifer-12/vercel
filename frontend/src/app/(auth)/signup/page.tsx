"use client";

import CustomInputGroup from "@/components/custom/input-group";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Field, FieldDescription, FieldGroup } from "@/components/ui/field";
import { useSignUp } from "@/hooks/use-auth";
import {
  EnvelopeIcon,
  IdentificationCardIcon,
  PasswordIcon,
} from "@phosphor-icons/react";

const SignupPage = () => {
  const { mutate, isPending, data } = useSignUp();
  return (
    <Card>
      <CardHeader>
        <CardTitle>Create an account</CardTitle>
        <CardDescription>
          Enter your information below to create your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            const form = e.currentTarget;
            const formData = new FormData(form);
            mutate(formData);
          }}
        >
          <FieldGroup>
            <Field>
              <CustomInputGroup
                id="name"
                placeholder="John Doe"
                text="Full Name"
                type="text"
                icon={IdentificationCardIcon}
                required={true}
                error={data?.error?.name?.[0] ?? null}
              />
            </Field>
            <Field>
              <CustomInputGroup
                icon={EnvelopeIcon}
                id="email"
                placeholder="john@gmail.com"
                text="Email"
                type="email"
                required={true}
                error={data?.error?.email?.[0] ?? null}
              />
            </Field>
            <Field>
              <CustomInputGroup
                icon={PasswordIcon}
                id="password"
                placeholder=""
                text="Password"
                type="password"
                required={true}
                error={data?.error?.password?.[0] ?? null}
              />
            </Field>
            <FieldGroup>
              <Field>
                <Button type="submit" disabled={isPending}>
                  Create Account
                </Button>
                <FieldDescription className="px-6 text-center">
                  Already have an account? <a href="/signin">Sign in</a>
                </FieldDescription>
              </Field>
            </FieldGroup>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
};

export default SignupPage;
