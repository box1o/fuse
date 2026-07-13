import { Badge } from "@/shared/components";
import { useListComputeServices } from "./hooks";

const Main = () => {
    const { data: services = [], error, isLoading } = useListComputeServices();

    if (isLoading) return <p className="p-6 text-muted-foreground">Loading compute services…</p>;
    if (error) return <p className="p-6 text-destructive">{error.message}</p>;

    return (
        <section className="mx-auto max-w-4xl p-6">
            <h1 className="text-2xl font-semibold">Compute</h1>
            <p className="mt-1 text-muted-foreground">Available compute services</p>
            <div className="mt-6 overflow-hidden rounded-lg border">
                {services.map((service) => (
                    <div key={service.id} className="flex items-center justify-between border-b p-4 last:border-b-0">
                        <span className="font-medium">{service.name}</span>
                        <Badge variant="secondary">{service.status}</Badge>
                    </div>
                ))}
                {services.length === 0 && <p className="p-4 text-muted-foreground">No compute services available.</p>}
            </div>
        </section>
    );
};

export default Main;
