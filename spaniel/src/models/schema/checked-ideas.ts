export type CheckedKeywordDto = {
  keyword: string;
  generatorVersions: number[];
  tokenlists: string[];
  runsCount: number;
};

export type CheckedKeyword = CheckedKeywordDto;

export type CheckedPasslistDto = {
  passlistUrl: string;
  name: string;
};

export type CheckedPasslist = CheckedPasslistDto;

export type CheckedIdeasDto = {
  checkedKeywords: CheckedKeyword[];
  checkedPasslists: CheckedPasslist[];
};

export type CheckedIdeas = CheckedIdeasDto;

export const mapCheckedIdeas = (dto: CheckedIdeasDto) => ({
  ...dto,
});
