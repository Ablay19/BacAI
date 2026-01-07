import { Subject } from "../types";
interface SubjectCardProps {
    subject: Subject;
    onSelect: (subject: Subject) => void;
    index: number;
}
export default function SubjectCard({ subject, onSelect, index, }: SubjectCardProps): import("react/jsx-runtime").JSX.Element;
interface SubjectsGridProps {
    subjects: Subject[];
    onSubjectSelect: (subject: Subject) => void;
}
export declare function SubjectsGrid({ subjects, onSubjectSelect }: SubjectsGridProps): import("react/jsx-runtime").JSX.Element;
export {};
//# sourceMappingURL=SubjectCard.d.ts.map