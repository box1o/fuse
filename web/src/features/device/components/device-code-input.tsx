import { InputOTP } from "@/shared/components";

interface DeviceCodeInputProps {
    value: string;
    disabled?: boolean;
    onChange: (value: string) => void;
}

const DeviceCodeInput = ({ value, disabled, onChange }: DeviceCodeInputProps) => {
    return (
        <InputOTP
            aria-label="Device authorization code"
            autoFocus
            containerClassName="justify-center"
            disabled={disabled}
            maxLength={8}
            onChange={onChange}
            value={value}
        >
            <InputOTP.Group>
                <InputOTP.Slot index={0} />
                <InputOTP.Slot index={1} />
                <InputOTP.Slot index={2} />
                <InputOTP.Slot index={3} />
            </InputOTP.Group>
            <InputOTP.Separator />
            <InputOTP.Group>
                <InputOTP.Slot index={4} />
                <InputOTP.Slot index={5} />
                <InputOTP.Slot index={6} />
                <InputOTP.Slot index={7} />
            </InputOTP.Group>
        </InputOTP>
    );
};

export { DeviceCodeInput };
