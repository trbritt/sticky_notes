# configure.ac

# Initialization
AC_INIT([gonotes], [1.0], [tristan.britt@mail.mcgill.ca])
AC_PROG_CC
AC_PROG_CXX
AC_PROG_RANLIB

AC_CONFIG_SRCDIR([src/main.rs])
AC_CONFIG_AUX_DIR(driver)

# Checks for programs
AC_CHECK_PROGS([GO], [go])
AC_CHECK_PROGS([RUST], [rustc])

# Checks for libraries
# (Add any additional library checks as needed)

# Checks for header files
# (Add any additional header file checks as needed)

# Checks for functions
# (Add any additional function checks as needed)

# Checks for Golang
if test -z "$GO"; then
  AC_MSG_ERROR([Golang compiler not found. Please install Go.])
fi

# Check for libgtk-4-dev
AC_CHECK_LIB([gtk-4], [gtk_init], [], [
  AC_MSG_ERROR([libgtk-4-dev not found. Please install the libgtk-4-dev package.])
])

# Checks for Rust
# (Add any additional Rust-related checks as needed)

# Configuration
srcdir=`dirname "$0"`
test -z "$srcdir" && srcdir=.

# Configuration
AC_CONFIG_FILES([
  Makefile
])

# Output information
AC_OUTPUT
