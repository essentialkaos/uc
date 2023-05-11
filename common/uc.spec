################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Simple utility for counting unique lines
Name:           uc
Version:        2.0.0
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/uc

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.19

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
Simple utility for counting unique lines.

################################################################################

%prep
%setup -q

%build
if [[ ! -d "%{name}/vendor" ]] ; then
  echo "This package requires vendored dependencies"
  exit 1
fi

pushd %{name}
  go build %{name}.go
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -dm 755 %{buildroot}%{_mandir}/man1

install -pm 755 %{name}/%{name} %{buildroot}%{_bindir}/

./%{name}/%{name} --generate-man > %{buildroot}%{_mandir}/man1/%{name}.1

%clean
rm -rf %{buildroot}

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_mandir}/man1/%{name}.1.*
%{_bindir}/%{name}

################################################################################

%changelog
* Wed May 03 2023 Anton Novojilov <andy@essentialkaos.com> - 2.0.0-0
- Better standard input handling

* Mon Mar 06 2023 Anton Novojilov <andy@essentialkaos.com> - 1.1.1-0
- Added verbose info output
- Dependencies update
- Code refactoring

* Thu Dec 01 2022 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-1
- Fixed build using sources from source.kaos.st

* Wed Aug 10 2022 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Minor UI improvements
- Fixed bug with parsing max number of unique lines
- Updated dependencies

* Tue Mar 29 2022 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Removed pkg.re usage
- Added module info
- Added Dependabot configuration

* Thu Oct 22 2020 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Added possibility to define -m/--max option as number with K and M
- Added man page generation

* Fri Jan 31 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.2-0
- Added option -m/--max for defining maximum unique lines to process
- Added option -d/--dist for printing data distribution

* Thu Jan 30 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- Initial build
