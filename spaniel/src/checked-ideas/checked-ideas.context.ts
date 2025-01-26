import {
  createContext,
  useContext
} from 'react';
import { CheckedIdeas } from '../models';

export type LoadCheckedIdeasFn = () => void;

export type CheckedIdeasContext = {
  loadIdeas: LoadCheckedIdeasFn;
  loading: boolean;
  checkedIdeas: CheckedIdeas;
};

export const checkedIdeasContext = createContext<CheckedIdeasContext>(null!);

export const useCheckedIdeasContext = () => useContext(checkedIdeasContext);
