import * as React from "react";
import { MinusIcon } from "lucide-react";
import { OTPInput as OTPInputPrimitive, OTPInputContext } from "input-otp";
import { cn } from "@/shared/utils";

type InputOTPProps = React.ComponentProps<typeof OTPInputPrimitive> & {
    containerClassName?: string;
};
type InputOTPGroupProps = React.ComponentProps<"div">;
type InputOTPSlotProps = React.ComponentProps<"div"> & {
    index: number;
};
type InputOTPSeparatorProps = React.ComponentProps<"div">;

const InputOTPRoot: React.FC<InputOTPProps> = ({ className, containerClassName, ...props }) => {
    return (
        <OTPInputPrimitive
            className={cn("disabled:cursor-not-allowed", className)}
            containerClassName={cn("flex items-center has-disabled:opacity-50", containerClassName)}
            data-slot="input-otp"
            spellCheck={false}
            {...props}
        />
    );
};

const InputOTPGroup: React.FC<InputOTPGroupProps> = ({ className, ...props }) => {
    return (
        <div
            className={cn(
                "flex items-center rounded-lg has-aria-invalid:border-destructive has-aria-invalid:ring-3 has-aria-invalid:ring-destructive/20 dark:has-aria-invalid:ring-destructive/40",
                className,
            )}
            data-slot="input-otp-group"
            {...props}
        />
    );
};

const InputOTPSlot: React.FC<InputOTPSlotProps> = ({ index, className, ...props }) => {
    const inputOTPContext = React.useContext(OTPInputContext);
    const { char, hasFakeCaret, isActive } = inputOTPContext?.slots[index] ?? {};

    return (
        <div
            className={cn(
                "relative flex size-11 items-center justify-center border-y border-r border-input font-mono text-lg transition-all outline-none first:rounded-l-lg first:border-l last:rounded-r-lg data-[active=true]:z-10 data-[active=true]:border-ring data-[active=true]:ring-3 data-[active=true]:ring-ring/50 aria-invalid:border-destructive data-[active=true]:aria-invalid:border-destructive data-[active=true]:aria-invalid:ring-destructive/20 dark:bg-input/30 dark:data-[active=true]:aria-invalid:ring-destructive/40",
                className,
            )}
            data-active={isActive}
            data-slot="input-otp-slot"
            {...props}
        >
            {char}
            {hasFakeCaret && (
                <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
                    <div className="h-5 w-px animate-caret-blink bg-foreground duration-1000" />
                </div>
            )}
        </div>
    );
};

const InputOTPSeparator: React.FC<InputOTPSeparatorProps> = ({ className, ...props }) => {
    return (
        <div
            className={cn("flex items-center px-1 text-muted-foreground [&_svg]:size-4", className)}
            data-slot="input-otp-separator"
            role="separator"
            {...props}
        >
            <MinusIcon />
        </div>
    );
};

const InputOTP = Object.assign(InputOTPRoot, {
    Group: InputOTPGroup,
    Slot: InputOTPSlot,
    Separator: InputOTPSeparator,
});

export { InputOTP };
