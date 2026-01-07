import { SolutionStep } from "../types";
interface SolutionDisplayProps {
    steps: SolutionStep[];
    finalAnswer: string;
    confidence: number;
    language: string;
    processingTime?: number;
}
export default function SolutionDisplay({ steps, finalAnswer, confidence, language, processingTime, }: SolutionDisplayProps): import("react/jsx-runtime").JSX.Element;
export {};
//# sourceMappingURL=SolutionDisplay.d.ts.map