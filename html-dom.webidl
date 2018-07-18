interface HTMLAnchorElement : HTMLElement {
           attribute DOMString target;
           attribute DOMString download;

           attribute DOMString rel;
           attribute DOMString rev;
  readonly attribute DOMTokenList relList;
           attribute DOMString hreflang;
           attribute DOMString type;

           attribute DOMString text;
};
HTMLAnchorElement implements URLUtils;
